package nav_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/nav"
)

func TestAngle(t *testing.T) {
	var tests = []struct {
		name    string
		in, out nav.Angle
	}{
		{
			name: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
