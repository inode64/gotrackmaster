package trackmaster

import (
	"math"

	gpx "github.com/twpayne/go-gpx"
)

// toRadians converts to radial coordinates.
func toRadians(x float64) float64 {
	return x / 180. * math.Pi
}

// toDegrees converts to degrees.
func toDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

// geoToCartesian converts geo coordinates to cartesian coordinates.
func geoToCartesian(coord gpx.WptType) (float64, float64, float64) {
	latRad := toRadians(coord.Lat)
	lonRad := toRadians(coord.Lon)
	r := earthRadius + coord.Ele
	x := r * math.Cos(latRad) * math.Cos(lonRad)
	y := r * math.Cos(latRad) * math.Sin(lonRad)
	z := r * math.Sin(latRad)
	return x, y, z
}

// cartesianToGeo converts cartesian coordinates to geo coordinates.
func cartesianToGeo(x, y, z float64) gpx.WptType {
	r := math.Sqrt(x*x + y*y + z*z)
	latRad := math.Asin(z / r)
	lonRad := math.Atan2(y, x)
	lat := toDegrees(latRad)
	lon := toDegrees(lonRad)
	alt := r - earthRadius

	return gpx.WptType{Lat: lat, Lon: lon, Ele: alt}
}

// mindpoint returns the midpoint between two coordinates.
func midpoint(coord1, coord2 gpx.WptType) gpx.WptType {
	x1, y1, z1 := geoToCartesian(coord1)
	x2, y2, z2 := geoToCartesian(coord2)
	xMid := (x1 + x2) / 2
	yMid := (y1 + y2) / 2
	zMid := (z1 + z2) / 2
	return cartesianToGeo(xMid, yMid, zMid)
}
