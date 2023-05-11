package trackmaster

import (
	"math"

	gpx "github.com/twpayne/go-gpx"
)

// Distance2D returns the 2D distance of two WptType.
func Distance2D(w, pt gpx.WptType) float64 {
	return distance(w.Lat, w.Lon, 0, pt.Lat, pt.Lon, 0, false, false)
}

// Distance3D returns the 3D distance of two WptType.
func Distance3D(w, pt gpx.WptType) float64 {
	return distance(w.Lat, w.Lon, w.Ele, pt.Lat, pt.Lon, pt.Ele, true, false)
}

// Distance returns the 2D or 3D distance of two WptType.
func distance(lat1, lon1, ele1, lat2, lon2, ele2 float64, threeD, haversine bool) float64 {
	absLat := math.Abs(lat1 - lat2)
	absLon := math.Abs(lon1 - lon2)
	if haversine || absLat > 0.2 || absLon > 0.2 {
		return HaversineDistance(lat1, lon1, lat2, lon2)
	}

	coef := math.Cos(toRadians(lat1))
	x := lat1 - lat2
	y := (lon1 - lon2) * coef

	distance2d := math.Sqrt(x*x+y*y) * oneDegree

	if !threeD || ele1 == ele2 {
		return distance2d
	}

	eleDiff := ele1 - ele2

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow(eleDiff, 2))
}

// HaversineDistance returns the haversine distance between two points.
//
// Implemented from http://www.movable-type.co.uk/scripts/latlong.html
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat1 - lat2)
	dLon := toRadians(lon1 - lon2)
	thisLat1 := toRadians(lat1)
	thisLat2 := toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c

	return d
}
