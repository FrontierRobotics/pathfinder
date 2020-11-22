package motor

type Movement int

const (
	Park        Movement = iota
	Forward              // Pair with Speed
	Reverse              // Pair with Speed
	RotateLeft           // Pair with Speed
	RotateRight          // Pair with Speed
)

type Speed int

const (
	Stop Speed = iota
	Slow
	Medium
	Full
)

type Command struct {
	M Movement
	S Speed
}
