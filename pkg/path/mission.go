package path

import (
	"time"

	"github.com/andycondon/pathfinder/pkg/nav"
	"github.com/golang/geo/s2"
)

type Waypoint struct {
	s2.LatLng
	Visited time.Time
}

type Mission struct {
	Waypoints []Waypoint
	idx       int
}

func NewMission(waypoints []s2.LatLng) *Mission {
	m := &Mission{
		Waypoints: make([]Waypoint, len(waypoints)),
	}
	for i, p := range waypoints {
		m.Waypoints[i] = Waypoint{LatLng: p}
	}
	return m
}

func (m Mission) VisitedWaypoints() []Waypoint {
	return m.Waypoints[:m.idx]
}

func (m *Mission) CurrentWaypoint() Waypoint {
	if !m.Complete() {
		return m.Waypoints[m.idx]
	}
	return Waypoint{}
}

func (m *Mission) Progress(position s2.LatLng) Waypoint {
	waypoint := m.CurrentWaypoint()
	distance := nav.DistanceBetweenPositions(position, waypoint.LatLng)
	if distance < 0.3*nav.Meters {
		m.Checkpoint()
	}
	return m.CurrentWaypoint()
}

func (m *Mission) Checkpoint() {
	if !m.Complete() {
		m.Waypoints[m.idx].Visited = time.Now()
		m.idx = m.idx + 1
	}
}

func (m *Mission) Complete() bool {
	return m.idx == len(m.Waypoints)
}
