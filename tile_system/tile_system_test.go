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

func TestLatLngToTileXY(t *testing.T) {
	cases := []struct {
		lat, lng float64
		zoom     uint
		x, y     int
	}{
		{
			-30.772727396604843, 22.463982922607443,
			3,
			4, 4,
		},

		{
			78.5555364791068, -166.85242332739256,
			3,
			0, 1,
		},
	}

	for _, c := range cases {
		x, y := LatLngToTileXY(c.lat, c.lng, c.zoom)
		if c.x != x || c.y != y {
			t.Errorf(
				"TestFail, expected {%d, %d}, got {%d, %d}",
				c.x, c.y,
				x, y,
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

func TestQuadKeyToTileXY(t *testing.T) {
	cases := []struct {
		quadKey string
		x, y    int
		zoom    uint
	}{
		{
			"213",
			3, 5, 3,
		},
		{
			"03130",
			14, 10, 5,
		},
	}

	for _, c := range cases {
		x, y, zoom := QuadKeyToTileXY(c.quadKey)
		if x != c.x || y != c.y || zoom != c.zoom {
			t.Errorf(
				"TestFail, expected {%d %d %u}, got {%d %d %u}",
				c.x, c.y, c.zoom,
				x, y, zoom,
			)
		}
	}
}
