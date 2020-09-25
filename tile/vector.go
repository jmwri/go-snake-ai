package tile

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
