package motor_test

import (
	"errors"
	"testing"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/andycondon/pathfinder/pkg/status"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	m1 := &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0}
	m2 := &motor.Motor{Addr: 0x02, Slow: 0x50, Med: 0xA0}
	bytesSent := make([]byte, 0)
	tx := func(w []byte, r []byte) error {
		for _, b := range w {
			bytesSent = append(bytesSent, b)
		}
		return nil
	}
	rs := func([]byte) (status.Reading, error) {
		return status.Reading{
			IR: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
		}, nil
	}
	driver := &motor.Driver{Tx: tx, ReadStatus: rs, Left: m1, Right: m2}

	t.Run("stop", func(t *testing.T) {
		_, _ = driver.Stop()

		assert.Equal(t, []byte{0x01, 0x01, 0x00, 0x02, 0x01, 0x00}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("forward full", func(t *testing.T) {
		_, _ = driver.Forward(motor.Full)

		assert.Equal(t, []byte{0x01, 0x01, 0xFF, 0x02, 0x01, 0xFF}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("reverse medium", func(t *testing.T) {
		_, _ = driver.Reverse(motor.Medium)

		assert.Equal(t, []byte{0x01, 0x00, 0xA0, 0x02, 0x00, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("rotate left medium", func(t *testing.T) {
		_, _ = driver.RotateLeft(motor.Medium)

		assert.Equal(t, []byte{0x01, 0x00, 0xA0, 0x02, 0x01, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("rotate right medium", func(t *testing.T) {
		_, _ = driver.RotateRight(motor.Medium)

		assert.Equal(t, []byte{0x01, 0x01, 0xA0, 0x02, 0x00, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("no error", func(t *testing.T) {
		_, err := driver.Forward(motor.Full)

		assert.NoError(t, err)

		bytesSent = make([]byte, 0)
	})

	t.Run("status returned", func(t *testing.T) {
		s, _ := driver.Forward(motor.Full)

		assert.Equal(t, status.Reading{
			IR: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
		}, s)

		bytesSent = make([]byte, 0)
	})

	t.Run("tx error", func(t *testing.T) {
		driver.Tx = func(w []byte, r []byte) error {
			return errors.New("kaboom")
		}

		_, err := driver.Forward(motor.Full)

		assert.Error(t, err)

		bytesSent = make([]byte, 0)
	})
}
