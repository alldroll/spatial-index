/*
 * inspired by https://msdn.microsoft.com/en-us/library/bb259689.aspx
 * https://developers.google.com/maps/documentation/javascript/examples/map-coordinates?hl=ru
 */

package tile_system

import (
	"bytes"
	"github.com/alldroll/spatial-index/geometry"
	"math"
)

const (
	tileSize    = 256
	earthRadius = 6378137
	minLat      = -85.05112878
	maxLat      = 85.05112878
	minLng      = -180.0
	maxLng      = 180.0
)

func GetScale(zoom uint) uint {
	return 1 << zoom
}

func GroundResolution(lat float64, zoom uint) float64 {
	latitute := clip(lat, minLat, maxLat)
	mapSize := GetScale(zoom) * tileSize
	return math.Cos(latitute*math.Pi/180) * 2 * math.Pi * earthRadius / float64(mapSize)
}

func LatLngToWorldCoordinate(lat, lng float64) *shape.Point {
	return Project(lat, lng)
}

func LatLngToPixelXY(lat, lng float64, zoom uint) (int, int) {
	scale := GetScale(zoom)
	worldCoordinate := LatLngToWorldCoordinate(lat, lng)
	x := math.Floor(worldCoordinate.GetX() * float64(scale))
	y := math.Floor(worldCoordinate.GetY() * float64(scale))
	return int(x), int(y)
}

func LatLngToTileXY(lat, lng float64, zoom uint) (int, int) {
	scale := GetScale(zoom)
	worldCoordinate := LatLngToWorldCoordinate(lat, lng)
	x := math.Floor(worldCoordinate.GetX() * float64(scale) / tileSize)
	y := math.Floor(worldCoordinate.GetY() * float64(scale) / tileSize)
	return int(x), int(y)
}

func TileXYToLatLng(x, y int, zoom uint) (float64, float64) {
	scale := GetScale(zoom)
	lng := float64(x)/float64(scale)*360 - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/float64(scale))))
	lat := latRad * 180.0 / math.Pi
	return lat, lng
}

func TileXYToQuadKey(x, y int, zoom uint) string {
	var buffer bytes.Buffer
	for i := zoom; i > 0; i-- {
		digit := '0'
		mask := 1 << (i - 1)
		if (x & mask) != 0 {
			digit++
		}

		if (y & mask) != 0 {
			digit += 2
		}

		buffer.WriteRune(digit)
	}

	return buffer.String()
}

func QuadKeyToTileXY(quadKey string) (int, int, uint) {
	zoom := uint(len(quadKey))
	x, y := 0, 0
	for i := zoom; i > 0; i-- {
		mask := 1 << (i - 1)
		switch quadKey[zoom-i] {
		case '0':
			break
		case '1':
			x |= mask
			break
		case '2':
			y |= mask
			break
		case '3':
			x |= mask
			y |= mask
			break
		default:
			panic("Invalid quadkey symbol")
		}
	}

	return x, y, zoom
}

func Project(lat, lng float64) *shape.Point {
	siny := math.Sin(lat * math.Pi / 180)
	siny = clip(siny, -0.9999, 0.9999)

	return shape.NewPoint(
		tileSize*(0.5+lng/360),
		tileSize*(0.5-math.Log((1+siny)/(1-siny))/(4*math.Pi)),
	)
}

func clip(n, min, max float64) float64 {
	return math.Min(math.Max(n, min), max)
}
