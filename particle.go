package main

import "github.com/go-gl/gl/v4.1-core/gl"

type Particle struct {
	drawable uint32

	x, y         float32
	mass         float32
	velocity     Vector
	acceleration Vector
	blackhole    bool
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

func newParticle(x, y, mass float32, blackhole bool) *Particle {
	points := make([]float32, len(square))
	copy(points, square)

	// Set the x and y position of the particle
	for i := 0; i < len(points); i += 3 {
		points[i] += x
		points[i+1] += y
	}

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
