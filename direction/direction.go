package direction

type Direction int

const (
	None Direction = iota
	Up
	Right
	Left
	Down
)

func Opposite(dir Direction) Direction {
	switch dir {
	case Up:
		return Down
	case Right:
		return Left
	case Down:
		return Up
	case Left:
		return Right
	}
	return None
}

func IsOpposite(d1 Direction, d2 Direction) bool {
	opp := Opposite(d1)
	return opp == d2
}
