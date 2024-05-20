package main

import "math"

type Vector struct {
	X, Y, Z float32
}

func (v Vector) add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v1 Vector) distance(v2 Vector) float32 {
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func (v Vector) normalize() Vector {
	magnitude := v.magnitude()
	return Vector{v.X / magnitude, v.Y / magnitude, v.Z / magnitude}
}

func (v Vector) magnitude() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector) mul(scalar float32) Vector {
	return Vector{v.X * scalar, v.Y * scalar, v.Z * scalar}
}
