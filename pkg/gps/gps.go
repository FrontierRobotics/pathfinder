package gps

import (
	"time"

	"github.com/andycondon/pathfinder/pkg/space"
)

type Reading struct {
	Fix      bool
	Time     time.Time
	Speed    space.Speed
	Position space.Position
}

func (r Reading) String() string {
	return "blah"
}

func FromGPRMC(sentence string) Reading {
	return Reading{}
}
