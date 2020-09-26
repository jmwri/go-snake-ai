package tile

import "go-snake-ai/direction"

func NewVector(x int, y int) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}

type Vector struct {
	X int
	Y int
}

func DirToVector(from *Vector, to *Vector) direction.Direction {
	if to.X == from.X && to.Y == from.Y-1 {
		return direction.Up
	}
	if to.X == from.X+1 && to.Y == from.Y {
		return direction.Right
	}
	if to.X == from.X && to.Y == from.Y+1 {
		return direction.Down
	}
	if to.X == from.X-1 && to.Y == from.Y {
		return direction.Left
	}
	return direction.None
}
