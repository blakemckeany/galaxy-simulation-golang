package main

type Quadtree struct {
	Bounds    Rect
	Nodes     [4]*Quadtree
	Max       int
	Particles []*Particle
}

type Rect struct {
	X, Y, W, H float32
}

func (r *Rect) Contains(p *Particle) bool {
	return p.x >= r.X && p.x <= r.X+r.W && p.y >= r.Y && p.y <= r.Y+r.H
}

func (r *Rect) Intersects(r2 Rect) bool {
	return r.X < r2.X+r2.W && r.X+r.W > r2.X && r.Y < r2.Y+r2.H && r.Y+r.H > r2.Y
}

const MaxParticles = 4

// Functions required
// 1. Insert
// 2. Divide
// 3. Query

// Create a new quadtree
func NewQuadtree(x, y, w, h float32) *Quadtree {
	return &Quadtree{
		Bounds:    Rect{x, y, w, h},
		Max:       MaxParticles,
		Particles: make([]*Particle, 0, MaxParticles),
	}
}

// Insert a particle into the quadtree
func (qt *Quadtree) Insert(p *Particle) {

	// If the particle is not in the bounds of the quadtree, return
	if !qt.Bounds.Contains(p) {
		return
	}

	// If the number of particles in the quadtree is less than the max allowed, add the particle to the quadtree
	if len(qt.Particles) < qt.Max {
		qt.Particles = append(qt.Particles, p)
		// Print the number of particles in the quadtree
		return
	}

	// If the number of particles in the quadtree is greater than the max allowed, divide the quadtree
	if len(qt.Particles) >= qt.Max {
		qt.Divide()
	}
}

// Divide the quadtree into 4 sub-quadtrees
func (qt *Quadtree) Divide() {
	x := qt.Bounds.X
	y := qt.Bounds.Y
	w := qt.Bounds.W / 2
	h := qt.Bounds.H / 2

	qt.Nodes[0] = NewQuadtree(x, y, w, h)
	qt.Nodes[1] = NewQuadtree(x+w, y, w, h)
	qt.Nodes[2] = NewQuadtree(x, y+h, w, h)
	qt.Nodes[3] = NewQuadtree(x+w, y+h, w, h)

	for _, p := range qt.Particles {
		for _, node := range qt.Nodes {
			node.Insert(p)
		}
	}

	qt.Particles = nil
}

// Get the quadtree node that the particle belongs to
func (qt *Quadtree) GetNode(p *Particle) *Quadtree {
	if !qt.Bounds.Contains(p) {
		return nil
	}

	for _, node := range qt.Nodes {
		if node.Bounds.Contains(p) {
			return node
		}
	}

	return nil
}

// Get neighboring nodes
func (qt *Quadtree) GetNeighbors() []*Quadtree {
	neighbors := make([]*Quadtree, 0, 4)

	for _, node := range qt.Nodes {
		if node != nil {
			neighbors = append(neighbors, node)
		}
	}

	return neighbors
}

// Query the quadtree for particles in a given range
func (qt *Quadtree) Query(r Rect) []*Particle {
	particles := make([]*Particle, 0, 10)

	if !qt.Bounds.Intersects(r) {
		return particles
	}

	for _, p := range qt.Particles {
		if r.Contains(p) {
			particles = append(particles, p)
		}
	}

	for _, node := range qt.Nodes {
		// Prevent null pointer exception
		if node == nil {
			continue
		}
		particles = append(particles, node.Query(r)...)
	}

	return particles
}
