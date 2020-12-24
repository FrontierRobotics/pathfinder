package motor_test

import (
	"errors"
	"testing"

	"github.com/andycondon/pathfinder/pkg/motor"
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
	driver := &motor.Driver{Tx: tx, Left: m1, Right: m2}

	t.Run("stop", func(t *testing.T) {
		_ = driver.D(motor.Command{M: motor.Park})

		assert.Equal(t, []byte{0x01, 0x01, 0x00, 0x02, 0x01, 0x00}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("forward full", func(t *testing.T) {
		_ = driver.D(motor.Command{M: motor.Forward, S: motor.Full})

		assert.Equal(t, []byte{0x01, 0x01, 0xFF, 0x02, 0x01, 0xFF}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("reverse medium", func(t *testing.T) {
		_ = driver.D(motor.Command{M: motor.Reverse, S: motor.Medium})

		assert.Equal(t, []byte{0x01, 0x00, 0xA0, 0x02, 0x00, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("rotate left medium", func(t *testing.T) {
		_ = driver.D(motor.Command{M: motor.RotateLeft, S: motor.Medium})

		assert.Equal(t, []byte{0x01, 0x00, 0xA0, 0x02, 0x01, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("rotate right medium", func(t *testing.T) {
		_ = driver.D(motor.Command{M: motor.RotateRight, S: motor.Medium})

		assert.Equal(t, []byte{0x01, 0x01, 0xA0, 0x02, 0x00, 0xA0}, bytesSent)

		bytesSent = make([]byte, 0)
	})

	t.Run("no error", func(t *testing.T) {
		err := driver.D(motor.Command{M: motor.Forward, S: motor.Full})

		assert.NoError(t, err)

		bytesSent = make([]byte, 0)
	})

	t.Run("tx error", func(t *testing.T) {
		driver.Tx = func(w []byte, r []byte) error {
			return errors.New("kaboom")
		}

		err := driver.D(motor.Command{M: motor.Forward, S: motor.Full})

		assert.Error(t, err)

		bytesSent = make([]byte, 0)
	})

	t.Run("bad command", func(t *testing.T) {
		driver.Tx = func(w []byte, r []byte) error {
			return errors.New("kaboom")
		}

		err := driver.D(motor.Command{M: motor.Movement(666)})

		assert.Error(t, err)

		bytesSent = make([]byte, 0)
	})
}
