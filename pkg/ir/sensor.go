package ir

import (
	"errors"
	"fmt"
)

type Level byte

func (l Level) Proximity() Proximity {
	if l < 0x10 {
		return ProximityClear
	}
	if l < 0x50 {
		return ProximityFar
	}
	return ProximityNear
}

type Proximity int

func (p Proximity) String() string {
	switch p {
	case ProximityNear:
		return "near"
	case ProximityFar:
		return "far"
	case ProximityClear:
		return "clear"
	default:
		return "unknown"
	}
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

type Sensor struct {
	Left    Level
	Forward Level
	Right   Level
}

func FromBytes(b []byte) (Sensor, error) {
	if len(b) != 3 {
		return Sensor{}, errors.New("sensor array requires 3 bytes")
	}
	return Sensor{
		Left:    Level(b[0]),
		Forward: Level(b[1]),
		Right:   Level(b[2]),
	}, nil
}

func (s *Sensor) Proximity() (l, f, r Proximity) {
	l = s.Left.Proximity()
	f = s.Forward.Proximity()
	r = s.Right.Proximity()
	return l, f, r
}

func (s *Sensor) String() string {
	l, f, r := s.Proximity()
	return fmt.Sprintf("L: %s,0x%02X F: %s,0x%02X R: %s,0x%02X", l, s.Left, f, s.Forward, r, s.Right)
}
