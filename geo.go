package trackmaster

import (
	"math"
)

// Distance2D returns the 2D distance of two WptType.
func (w *WptType) Distance2D(pt WptType) float64 {
	return distance(w.Lat, w.Lon, 0, pt.Lat, pt.Lon, 0, false, false)
}

// Distance3D returns the 3D distance of two WptType.
func (w *WptType) Distance3D(pt WptType) float64 {
	return distance(w.Lat, w.Lon, w.Ele, pt.Lat, pt.Lon, pt.Ele, true, false)
}

// toRadians converts to radial coordinates.
func toRadians(x float64) float64 {
	return x / 180. * math.Pi
}

func toDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func geoToCartesian(coord WptType) (float64, float64, float64) {
	latRad := toRadians(coord.Lat)
	lonRad := toRadians(coord.Lon)

	r := earthRadius + coord.Ele

	x := r * math.Cos(latRad) * math.Cos(lonRad)
	y := r * math.Cos(latRad) * math.Sin(lonRad)
	z := r * math.Sin(latRad)

	return x, y, z
}

func cartesianToGeo(x, y, z float64) WptType {
	r := math.Sqrt(x*x + y*y + z*z)
	latRad := math.Asin(z / r)
	lonRad := math.Atan2(y, x)

	lat := toDegrees(latRad)
	lon := toDegrees(lonRad)
	alt := r - earthRadius

	return WptType{Lat: lat, Lon: lon, Ele: alt}
}
