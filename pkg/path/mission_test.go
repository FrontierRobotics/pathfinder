package path

import (
	"testing"

	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestMission(t *testing.T) {
	var (
		s  = s2.LatLngFromDegrees(41.185579119692925, -104.80847684264579)
		m1 = s2.LatLngFromDegrees(41.18553067368046, -104.80847416049401)
		w1 = s2.LatLngFromDegrees(41.185493329876394, -104.80846343169397)
		m2 = s2.LatLngFromDegrees(41.18550645031859, -104.80825556017193)
		w2 = s2.LatLngFromDegrees(41.18551556476105, -104.80804562107231)
	)

	t.Run("start", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)

		w := m.CurrentWaypoint()
		assert.Equal(t, w1, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 0)
		assert.False(t, m.Complete())
	})

	t.Run("progress to midpoint 1", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)

		w := m.CurrentWaypoint()
		assert.Equal(t, w1, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 0)
		assert.False(t, m.Complete())
	})

	t.Run("progress to waypoint 1", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)
		m.Progress(w1)

		w := m.CurrentWaypoint()
		assert.Equal(t, w2, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 1)
		assert.False(t, m.Complete())
	})

	t.Run("revisits waypoint 1", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)
		m.Progress(w1)
		m.Progress(w1)

		w := m.CurrentWaypoint()
		assert.Equal(t, w2, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 1)
		assert.False(t, m.Complete())
	})

	t.Run("progress to midpoint 2", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)
		m.Progress(w1)
		m.Progress(m2)

		w := m.CurrentWaypoint()
		assert.Equal(t, w2, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 1)
		assert.False(t, m.Complete())
	})

	t.Run("progress to waypoint 2", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)
		m.Progress(w1)
		m.Progress(m2)
		m.Progress(w2)

		w := m.CurrentWaypoint()
		assert.Equal(t, s2.LatLng{}, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 2)
		assert.True(t, m.Complete())
	})

	t.Run("revisits waypoint 2", func(t *testing.T) {
		m := NewMission([]s2.LatLng{w1, w2})

		m.Progress(s)
		m.Progress(m1)
		m.Progress(w1)
		m.Progress(m2)
		m.Progress(w2)
		m.Progress(w2)

		w := m.CurrentWaypoint()
		assert.Equal(t, s2.LatLng{}, w.LatLng)
		assert.Len(t, m.VisitedWaypoints(), 2)
		assert.True(t, m.Complete())
	})
}
