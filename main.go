package main

import (
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"strconv"

	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 900
	height = 900

	particleCount     = 25000
	particleSize      = 0.0015
	grav              = 0.01
	particleMass      = 1.0
	blackHoleMass     = 999999999999
	softeningFactor   = 0.9
	dampeningFactor   = 0.9983
	initialSpinFactor = 0.1
	queryRange        = 15

	// Vertex shader, GLSL
	vertexShaderSource = `
        #version 410
        in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	// Fragment shader, GLSL
	fragmentShaderSource = `
        #version 410
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 1, 1, 0.6);
        }
    ` + "\x00"
)

type Particle struct {
	drawable uint32

	x, y         float32
	mass         float32
	velocity     Vector
	acceleration Vector
	blackhole    bool
}

var (
	// Slice of vertices for a square
	square = []float32{

		-particleSize, particleSize, 0, // top
		-particleSize, -particleSize, 0, // left
		particleSize, -particleSize, 0, // right

		-particleSize, particleSize, 0, // top
		particleSize, particleSize, 0, // left
		particleSize, -particleSize, 0, // right
	}

	quadTree = NewQuadtree(-1, -1, 2, 2)

	frames      = float32(0)
	elapsedTime = float32(0)

	lastTime = float32(0.0)
)

func main() {

	rand.Seed(time.Now().UnixNano())

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	vao := makeParticles()

	for !window.ShouldClose() {
		// TODO
		draw(vao, window, program)
	}
}

func (p *Particle) draw() {
	// Update the vertex data based on the current position
	points := make([]float32, len(square))
	copy(points, square)

	for i := 0; i < len(points); i += 3 {
		points[i] += p.x
		points[i+1] += p.y
	}

	// Update the vertex data in the VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, p.drawable)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// Render the particle
	gl.BindVertexArray(p.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func draw(particles []*Particle, window *glfw.Window, prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	for _, p := range particles {
		p.draw()
	}

	currentTime := float32(glfw.GetTime())
	dt := float32(currentTime - lastTime)
	lastTime = currentTime

	elapsedTime += dt
	frames++

	if elapsedTime >= 1.0 {
		fps := frames / elapsedTime
		window.SetTitle("Galaxy Simulation - FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 64) + " - Particles: " + strconv.Itoa(particleCount))
		frames = 0
		elapsedTime = 0
	}

	updateParticles(particles, dt)

	// updateParticles(particles)

	glfw.PollEvents()
	window.SwapBuffers()
}

func updateParticles(particles []*Particle, dt float32) {

	// Quad tree
	quadTree = NewQuadtree(-1, -1, 2, 2)

	for _, p := range particles {
		quadTree.Insert(p)
	}

	// Calculate the forces between the particles
	for _, p := range particles {

		p.acceleration = Vector{0, 0, 0} // Reset acceleration

		queryRange := Rect{p.x - queryRange/2, p.y - queryRange/2, queryRange, queryRange}
		nearbyParticles := quadTree.Query(queryRange)

		// node := quadTree.GetNode(p)

		// // If not null, get the particles in the node and its neighbours
		// if node == nil {
		// 	continue
		// }
		// particlesInNode := node.Particles

		// neighbourNodes := node.GetNeighbors()

		// for _, neighbour := range neighbourNodes {
		// 	particlesInNode = append(particlesInNode, neighbour.Particles...)
		// }

		// if node != nil {
		for _, p2 := range nearbyParticles {
			if p == p2 {
				continue
			}

			// Calculate the force between the particles
			force := calculateForce(p, p2)

			// Accumulate the force to the particle's acceleration
			p.acceleration = p.acceleration.add(force.mul(0.8 / p.mass))
		}
		// }
	}

	// Update the position and velocity of the particles
	for _, p := range particles {
		// Update velocity based on acceleration
		p.velocity = p.velocity.add(p.acceleration.mul(dt))

		// Dampen the velocity
		p.velocity = p.velocity.mul(dampeningFactor)

		// Update position based on velocity

		p.x += p.velocity.X * dt
		p.y += p.velocity.Y * dt

		// Wrap the particles around the screen
		if p.x < -1 {
			p.x = 1
		}
		if p.x > 1 {
			p.x = -1
		}
		if p.y < -1 {
			p.y = 1
		}
		if p.y > 1 {
			p.y = -1
		}

	}
}

func calculateForce(p1, p2 *Particle) Vector {
	// Calculate the distance between the particles
	distance := Vector{p2.x - p1.x, p2.y - p1.y, 0}
	distanceMagnitude := distance.magnitude()

	// Apply a softening factor to prevent excessive force at close distances
	softening := float32(softeningFactor)
	distanceMagnitude = float32(math.Sqrt(float64(distanceMagnitude*distanceMagnitude + softening*softening)))

	// Calculate the force magnitude using Newton's law of universal gravitation
	forceMagnitude := grav * p1.mass * p2.mass / (distanceMagnitude * distanceMagnitude)

	// Calculate the force vector by normalizing the distance vector and scaling it by the force magnitude
	forceVector := distance.normalize().mul(forceMagnitude)

	return forceVector
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func makeParticles() []*Particle {
	particles := make([]*Particle, particleCount)
	for x := 0; x < particleCount; x++ {
		// Random X, Random Y between -1.0 and 1.0
		xPosition := rand.Float32()*2 - 1
		yPosition := rand.Float32()*2 - 1

		// Float32 mass of 100.0
		mass := float32(particleMass)
		blackhole := false

		if x == 0 {
			xPosition = 0
			yPosition = 0
			mass = blackHoleMass
			blackhole = true
		}

		p := newParticle(xPosition, yPosition, mass, blackhole)
		particles[x] = p
	}

	return particles
}

func newParticle(x, y, mass float32, blackhole bool) *Particle {
	points := make([]float32, len(square))
	copy(points, square)

	// Set the x and y position of the particle
	for i := 0; i < len(points); i += 3 {
		points[i] += x
		points[i+1] += y
	}

	// Random velocity

	// Calculate the initial velocity of the particle
	// Based on perpindicular vector to the center of the screen
	// This will give the particle a circular orbit

	center := Vector{0, 0, 0}
	position := Vector{x, y, 0}
	distance := position.distance(center)

	// Calculate the velocity vector
	velocity := Vector{0, 0, 0}
	if distance != 0 {
		velocity = Vector{-(y - center.Y) / distance, (x - center.X) / distance, 0}.mul(initialSpinFactor)
	}

	return &Particle{
		drawable: makeVao(points),

		x:         x,
		y:         y,
		velocity:  velocity,
		mass:      mass,
		blackhole: blackhole,
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	// Enable blending
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Galaxy Simulation", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}
