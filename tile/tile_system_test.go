package tile_system

import "testing"

func TestMercatorProjection(t *testing.T) {
	cases := []struct {
		lat, lng float64
		x, y     float64
	}{
		{
			41.850, -87.650,
			65.67111111111113, 95.17492654697409,
		},
	}

	for _, c := range cases {
		point := Project(c.lat, c.lng)

		if !point.EqualXY(c.x, c.y) {
			t.Errorf(
				"TestFail, expected {x: %g, y: %g}, got {x: %g, y: %g}",
				c.x, c.y,
				point.GetX(), point.GetY(),
			)
		}
	}
}

func TestTileXYToQuadKey(t *testing.T) {
	cases := []struct {
		x, y    int
		zoom    uint
		quadKey string
	}{
		{
			3, 5, 3,
			"213",
		},
		{
			14, 10, 5,
			"03130",
		},
	}

	for _, c := range cases {
		actual := TileXYToQuadKey(c.x, c.y, c.zoom)
		if actual != c.quadKey {
			t.Errorf(
				"TestFail, expected {%s}, got {%s}",
				c.quadKey,
				actual,
			)
		}
	}
}
