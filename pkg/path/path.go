package path

import "github.com/andycondon/pathfinder/pkg/ir"

type Path string

const (
	Clear   Path = "clear"
	Forward      = "forward"
	Left         = "left"
	Right        = "right"
	Blocked      = "blocked"
)

func Find(irSensor ir.Sensor) Path {
	l, f, r := irSensor.Proximity()

	if l.IsClear() && f.IsClear() && r.IsClear() {
		return Clear
	}
	if f < ir.ProximityNear && (l.IsFar() || r.IsFar()) {
		return Forward
	}
	return Blocked
}
