package path

import "github.com/andycondon/pathfinder/pkg/ir"

type Path string

const (
	Forward Path = "forward"
	Left         = "left"
	Right        = "right"
	Blocked      = "blocked"
)

func Find(r ir.Reading) Path {
	if r.AllClear() {
		return Forward
	}
	if r.F < ir.ProximityNear && (r.L.IsFar() || r.R.IsFar()) {
		return Forward
	}
	return Blocked
}
