package ir_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/stretchr/testify/assert"
)

func TestSensorArray_Reading(t *testing.T) {
	tests := []struct {
		name    string
		reading ir.Reading
		bytes   []byte
		hasErr  bool
	}{
		{
			name:    "all clear",
			bytes:   []byte{0x10, 0x20, 0x30},
			reading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
		},
		{
			name:    "all far",
			bytes:   []byte{0x50, 0x60, 0x70},
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
		},
		{
			name:    "all near",
			bytes:   []byte{0x51, 0x61, 0x71},
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
		},
		{
			name:   "bad input",
			bytes:  []byte{0x50, 0x60},
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sensorArray := &ir.SensorArray{
				Left:    ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
				Forward: ir.Sensor{ClearUpperBound: 0x20, FarUpperBound: 0x60},
				Right:   ir.Sensor{ClearUpperBound: 0x30, FarUpperBound: 0x70},
			}
			reading, err := sensorArray.Reading(tc.bytes)
			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.reading, reading)
		})
	}
}

func TestProximity_String(t *testing.T) {
	tests := []struct {
		name    string
		reading ir.Reading
		output  string
	}{
		{
			name:    "all clear",
			reading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
			output:  "L: clear F: clear R: clear",
		},
		{
			name:    "all far",
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
			output:  "L: far F: far R: far",
		},
		{
			name:    "all near",
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
			output:  "L: near F: near R: near",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.reading.String())
		})
	}
}

func TestReading_AllClear(t *testing.T) {
	tests := []struct {
		reading ir.Reading
		isClear bool
	}{
		{
			reading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
			isClear: true,
		},
		{
			reading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityNear},
			isClear: false,
		},
		{
			reading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityNear, R: ir.ProximityNear},
			isClear: false,
		},
		{
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
			isClear: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.reading.String(), func(t *testing.T) {
			assert.Equal(t, tc.isClear, tc.reading.AllClear())
		})
	}
}

func TestReading_AllFar(t *testing.T) {
	tests := []struct {
		reading ir.Reading
		isFar   bool
	}{
		{
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
			isFar:   true,
		},
		{
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityNear},
			isFar:   false,
		},
		{
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityNear, R: ir.ProximityNear},
			isFar:   false,
		},
		{
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
			isFar:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.reading.String(), func(t *testing.T) {
			assert.Equal(t, tc.isFar, tc.reading.AllFar())
		})
	}
}

func TestReading_AllNear(t *testing.T) {
	tests := []struct {
		reading ir.Reading
		isNear  bool
	}{
		{
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
			isNear:  true,
		},
		{
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityFar},
			isNear:  false,
		},
		{
			reading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityFar, R: ir.ProximityFar},
			isNear:  false,
		},
		{
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
			isNear:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.reading.String(), func(t *testing.T) {
			assert.Equal(t, tc.isNear, tc.reading.AllNear())
		})
	}
}
