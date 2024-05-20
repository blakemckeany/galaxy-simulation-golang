package main

type Quadtree struct {
	Bounds Rect
	Nodes  [4]*Quadtree
	Max    int
}

type Rect struct {
	X, Y, W, H float32
}

const MaxObjects = 5
