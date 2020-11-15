package ir

import (
	"errors"
	"fmt"
)

type Sensor struct {
	ClearUpperBound, FarUpperBound byte
}

func (s Sensor) Proximity(l byte) Proximity {
	if l <= s.ClearUpperBound {
		return ProximityClear
	}
	if l <= s.FarUpperBound {
		return ProximityFar
	}
	return ProximityNear
}

type Proximity int

func (p Proximity) String() string {
	if p.IsClear() {
		return "clear"
	}
	if p.IsFar() {
		return "far"
	}
	return "near"
}

const (
	ProximityClear Proximity = iota
	ProximityFar
	ProximityNear
)

func (p Proximity) IsClear() bool {
	return p == ProximityClear
}

func (p Proximity) IsFar() bool {
	return p == ProximityFar
}

func (p Proximity) IsNear() bool {
	return p == ProximityNear
}

type Reading struct {
	L, F, R Proximity
}

func (r Reading) AllClear() bool {
	return r.L.IsClear() && r.F.IsClear() && r.R.IsClear()
}

func (r Reading) AllFar() bool {
	return r.L.IsFar() && r.F.IsFar() && r.R.IsFar()
}

func (r Reading) AllNear() bool {
	return r.L.IsNear() && r.F.IsNear() && r.R.IsNear()
}

func (r Reading) String() string {
	return fmt.Sprintf("L: %s F: %s R: %s", r.L.String(), r.F.String(), r.R.String())
}

type SensorArray struct {
	Left    Sensor
	Forward Sensor
	Right   Sensor
}

func (s *SensorArray) Reading(b []byte) (Reading, error) {
	if len(b) != 3 {
		return Reading{}, errors.New("sensor array requires 3 bytes")
	}
	return Reading{
		L: s.Left.Proximity(b[0]),
		F: s.Forward.Proximity(b[1]),
		R: s.Right.Proximity(b[2]),
	}, nil
}
