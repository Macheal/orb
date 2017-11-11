package project

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestToPlanar(t *testing.T) {
	for _, g := range orb.AllGeometries {
		// should not panic with unsupported type
		ToPlanar(g, Mercator)
	}
}

func TestToGeo(t *testing.T) {
	for _, g := range orb.AllGeometries {
		// should not panic with unsupported type
		ToGeo(g, Mercator)
	}
}
