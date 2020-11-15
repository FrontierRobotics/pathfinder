package path_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/path"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		irReading ir.Reading
		path      path.Path
	}{
		{
			name:      "move forward if all clear",
			irReading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
			path:      path.Forward,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := path.Find(tc.irReading)

			assert.Equal(t, tc.path, path)
		})
	}
}
