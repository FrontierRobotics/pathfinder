package motor_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/stretchr/testify/assert"
)

func TestMotor_Forward(t *testing.T) {
	tests := []struct {
		name    string
		speed   motor.Speed
		command []byte
	}{
		{
			name:    "full",
			speed:   motor.Full,
			command: []byte{0x01, 0x01, 0xFF},
		},
		{
			name:    "medium",
			speed:   motor.Medium,
			command: []byte{0x01, 0x01, 0xA0},
		},
		{
			name:    "slow",
			speed:   motor.Slow,
			command: []byte{0x01, 0x01, 0x50},
		},
		{
			name:    "stop",
			speed:   motor.Stop,
			command: []byte{0x01, 0x01, 0x00},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0}

			assert.Equal(t, tc.command, m.Forward(tc.speed))
		})
	}
}

func TestMotor_Reverse(t *testing.T) {
	tests := []struct {
		name    string
		speed   motor.Speed
		command []byte
	}{
		{
			name:    "full",
			speed:   motor.Full,
			command: []byte{0x01, 0x00, 0xFF},
		},
		{
			name:    "medium",
			speed:   motor.Medium,
			command: []byte{0x01, 0x00, 0xA0},
		},
		{
			name:    "slow",
			speed:   motor.Slow,
			command: []byte{0x01, 0x00, 0x50},
		},
		{
			name:    "stop",
			speed:   motor.Stop,
			command: []byte{0x01, 0x00, 0x00},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0}

			assert.Equal(t, tc.command, m.Reverse(tc.speed))
		})
	}
}

func TestMotor_Stop(t *testing.T) {
	m := &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0}

	assert.Equal(t, []byte{0x01, 0x01, 0x00}, m.Stop())
}
