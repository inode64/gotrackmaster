package trackmaster

import (
	"fmt"
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
	coefficient := math.Cos(toRadians(lat1))
	x := lat1 - lat2
	y := (lon1 - lon2) * coefficient
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

func HaversineDistanceTrkPt(pointA, pointB gpx.WptType) float64 {
	return HaversineDistance(pointA.Lat, pointA.Lon, pointB.Lat, pointB.Lon)
}

// Gaussian smooths the positions of a GPX file using a Gaussian filter.
func SmoothGaussian(g gpx.GPX, windowSize int, sigma float64) {
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			gaussianFilterPositions(*TrkSegType, windowSize, sigma)
		}
	}
}

func findNextCloserPoint(ts gpx.TrkSegType, start, max int, maxDistance, maxElevation float64) (int, float64) {
	lastPoint := -1
	minDistance := math.MaxFloat64
	for i := start + 1; i < MinInt(start+max, len(ts.TrkPt)); i++ {
		distance := HaversineDistanceTrkPt(*ts.TrkPt[start], *ts.TrkPt[i])
		elevation := ElevationAbs(*ts.TrkPt[start], *ts.TrkPt[i])

		if distance < minDistance && distance < maxDistance && elevation <= maxElevation {
			minDistance = distance
			lastPoint = i
		}
	}

	if lastPoint == -1 {
		return -1, math.MaxFloat64
	}

	return lastPoint, minDistance
}

func gaussianFilterPositions(position gpx.TrkSegType, windowSize int, sigma float64) {
	smoothed := make([]gpx.WptType, len(position.TrkPt))
	for i := 0; i < len(position.TrkPt); i++ {
		weights := make([]float64, len(position.TrkPt))
		sumWeights := 0.0
		normLat, normLon := 0.0, 0.0
		for j := -windowSize; j < windowSize; j++ {
			if i-windowSize/2+j < 0 || i+windowSize/2+j >= len(position.TrkPt) {
				continue
			}
			weights[i-windowSize/2+j] = Gaussian(float64(j-windowSize/2), sigma)
			sumWeights += weights[i-windowSize/2+j]

			normLat += weights[i-windowSize/2+j] * position.TrkPt[i-windowSize/2+j].Lat
			normLon += weights[i-windowSize/2+j] * position.TrkPt[i-windowSize/2+j].Lon
		}
		smoothed[i].Lat = normLat / sumWeights
		smoothed[i].Lon = normLon / sumWeights
	}
	for i := 0; i < len(position.TrkPt); i++ {
		if i >= len(position.TrkPt) {
			continue
		}
		position.TrkPt[i].Lat = smoothed[i].Lat
		position.TrkPt[i].Lon = smoothed[i].Lon
	}
}

// remove points when accuracy is too low in first point.
func RemoveFirstNoise(g gpx.GPX, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			var dst []*gpx.WptType
			// not enough points
			if len(TrkSegType.TrkPt) < MinSegmentLength {
				continue
			}
			for i := 0; i < 11; i++ {
				nextDistance := HaversineDistanceTrkPt(*TrkSegType.TrkPt[i], *TrkSegType.TrkPt[i+1])
				closerPoint, closerDistance := findNextCloserPoint(*TrkSegType, i, 5, 8, 0)
				if nextDistance > closerDistance {
					point := GPXElementInfo{
						WptType:      *TrkSegType.TrkPt[i],
						WptTypeNo:    i,
						TrkSegTypeNo: TrkSegTypeNo,
						TrkTypeNo:    TrkTypeNo,
					}
					result = append(result, point)
					if fix {
						dst = append(dst, TrkSegType.TrkPt[i])
						if closerPoint >= 10 {
							dst = append(dst, TrkSegType.TrkPt[closerPoint:]...)
						} else {
							dst = append(dst, TrkSegType.TrkPt[closerPoint])
						}
					}
					i = closerPoint
				} else if fix {
					if i >= 10 {
						dst = append(dst, TrkSegType.TrkPt[i:]...)
					} else {
						dst = append(dst, TrkSegType.TrkPt[i])
					}
				}
			}
			if fix && len(dst) > 0 {
				g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt = dst
			}
		}
	}
	return result
}

func RemoveNoise(g gpx.GPX, maxDistance, maxElevation float64, maxPoints int, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			var dst []*gpx.WptType
			for wptTypeNo := 0; wptTypeNo < len(TrkSegType.TrkPt)-1; wptTypeNo++ {
				nextDistance := HaversineDistanceTrkPt(*TrkSegType.TrkPt[wptTypeNo], *TrkSegType.TrkPt[wptTypeNo+1])
				closerPoint, closerDistance := findNextCloserPoint(*TrkSegType, wptTypeNo, maxPoints, maxDistance, maxElevation)
				if nextDistance > closerDistance {
					point := GPXElementInfo{
						WptType:      *TrkSegType.TrkPt[wptTypeNo],
						WptTypeNo:    wptTypeNo,
						TrkSegTypeNo: TrkSegTypeNo,
						TrkTypeNo:    TrkTypeNo,
					}
					result = append(result, point)
					dst = append(dst, TrkSegType.TrkPt[wptTypeNo])
					dst = append(dst, TrkSegType.TrkPt[closerPoint])
					wptTypeNo = closerPoint
				} else {
					dst = append(dst, TrkSegType.TrkPt[wptTypeNo])
				}
			}
			if fix && len(dst) > 0 {
				g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt = dst
			}
		}
	}
	return result
}

func RemoveStops(g gpx.GPX, minSeconds, maxDistance, maxElevation float64, minPoints int, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	var distance float64
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			var dst []*gpx.WptType
			firstPoint := -1
			var numPoints, point int
			for wptTypeNo := 0; wptTypeNo < len(TrkSegType.TrkPt)-1; wptTypeNo++ {
				if firstPoint == -1 {
					point = wptTypeNo
				} else {
					point = firstPoint
				}
				if TrkSegType.TrkPt[point].Ele == 0 {
					TrkSegType.TrkPt[point].Ele = TrkSegType.TrkPt[wptTypeNo+1].Ele
				}
				distance = HaversineDistanceTrkPt(*TrkSegType.TrkPt[point], *TrkSegType.TrkPt[wptTypeNo+1])
				elevation := ElevationAbs(*TrkSegType.TrkPt[point], *TrkSegType.TrkPt[wptTypeNo+1])
				if distance <= maxDistance && elevation <= maxElevation {
					if firstPoint == -1 {
						firstPoint = wptTypeNo
					}
					numPoints++
				} else {
					seconds := TimeDiff(*TrkSegType.TrkPt[point], *TrkSegType.TrkPt[wptTypeNo])
					if numPoints > minPoints && seconds > minSeconds {
						distance = HaversineDistanceTrkPt(*TrkSegType.TrkPt[firstPoint], *TrkSegType.TrkPt[wptTypeNo])
						elevation = ElevationAbs(*TrkSegType.TrkPt[firstPoint], *TrkSegType.TrkPt[wptTypeNo])
						point := GPXElementInfo{
							WptType:      *TrkSegType.TrkPt[firstPoint],
							WptTypeNo:    firstPoint,
							TrkSegTypeNo: TrkSegTypeNo,
							TrkTypeNo:    TrkTypeNo,
							Count:        numPoints,
							Length:       distance,
							Elevation:    elevation,
							Duration:     seconds,
						}
						result = append(result, point)
						if numPoints > minPoints && seconds > minSeconds {
							dst = append(dst, TrkSegType.TrkPt[firstPoint])
						} else {
							dst = append(dst, TrkSegType.TrkPt[firstPoint:wptTypeNo+1]...)
						}
						// for remove close points
						if minPoints != 0 {
							dst = append(dst, TrkSegType.TrkPt[wptTypeNo])
						}
					} else {
						if firstPoint == -1 {
							dst = append(dst, TrkSegType.TrkPt[wptTypeNo])
						} else {
							dst = append(dst, TrkSegType.TrkPt[firstPoint:wptTypeNo+1]...)
						}
					}
					firstPoint, numPoints = -1, 0
				}
			}
			if fix {
				if numPoints == 0 {
					if len(TrkSegType.TrkPt) != 0 {
						dst = append(dst, TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1])
					}
				} else {
					dst = append(dst, TrkSegType.TrkPt[firstPoint:]...)
					distance = HaversineDistanceTrkPt(*TrkSegType.TrkPt[firstPoint], *TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1])
					elevation := ElevationAbs(*TrkSegType.TrkPt[firstPoint], *TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1])
					seconds := TimeDiff(*TrkSegType.TrkPt[firstPoint], *TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1])
					point := GPXElementInfo{
						WptType:      *TrkSegType.TrkPt[firstPoint],
						WptTypeNo:    firstPoint,
						TrkSegTypeNo: TrkSegTypeNo,
						TrkTypeNo:    TrkTypeNo,
						Count:        numPoints,
						Length:       distance,
						Elevation:    elevation,
						Duration:     seconds,
					}
					result = append(result, point)
				}
				g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt = dst
			}
		}
	}

	return result
}

// function to check if the two segments p1q1 and p2q2 intersect.
func doIntersect(p1, q1, p2, q2 gpx.WptType) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 && o1 != 0 && o2 != 0 && o3 != 0 && o4 != 0 {
		return true
	}

	return false
}

// function to find orientation of ordered triplet (p, q, r).
func orientation(p, q, r gpx.WptType) int {
	val := (q.Lon-p.Lon)*(r.Lat-q.Lat) - (q.Lat-p.Lat)*(r.Lon-q.Lon)

	if val == 0 {
		// returns Colinear points
		return 0
	}

	if val > 0 {
		// returns Clockwise points
		return 1
	}

	// returns Counterclockwise
	return 2
}

// CheckIntersecting - check intersecting segments.
func RemoveIntersections(g gpx.GPX, max int, fix bool) []GPXElementInfo {
	var result []GPXElementInfo

	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo := 0; wptTypeNo < len(TrkSegType.TrkPt)-1; wptTypeNo++ {
				lastPoint := -1
				for j := wptTypeNo + 2; j < MinInt(wptTypeNo+max, len(TrkSegType.TrkPt)-1); j++ {
					if doIntersect(*TrkSegType.TrkPt[wptTypeNo], *TrkSegType.TrkPt[wptTypeNo+1], *TrkSegType.TrkPt[j], *TrkSegType.TrkPt[j+1]) {
						point := GPXElementInfo{
							WptType:      *TrkSegType.TrkPt[wptTypeNo],
							WptTypeNo:    wptTypeNo,
							TrkSegTypeNo: TrkSegTypeNo,
							TrkTypeNo:    TrkTypeNo,
						}
						result = append(result, point)
						lastPoint = j + 1
						break
					}
				}
				if lastPoint != -1 {
					if fix {
						TrkSegType.TrkPt = append(TrkSegType.TrkPt[:wptTypeNo+1], TrkSegType.TrkPt[lastPoint:]...)
					}
					wptTypeNo = lastPoint - 1
				}
			}
		}
	}
	return result
}

// Get bounds of GPX.
func GetBounds(g gpx.GPX) gpx.BoundsType {
	var result gpx.BoundsType
	result.MinLat = 90
	result.MaxLat = -90
	result.MinLon = 180
	result.MaxLon = -180
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for _, WptType := range TrkSegType.TrkPt {
				if WptType.Lat < result.MinLat {
					result.MinLat = WptType.Lat
				}
				if WptType.Lat > result.MaxLat {
					result.MaxLat = WptType.Lat
				}
				if WptType.Lon < result.MinLon {
					result.MinLon = WptType.Lon
				}
				if WptType.Lon > result.MaxLon {
					result.MaxLon = WptType.Lon
				}
			}
		}
	}
	return result
}

func IsBoundsValid(bounds gpx.BoundsType) bool {
	return bounds.MaxLat != 90 || bounds.MinLat != -90 || bounds.MaxLon != 180 || bounds.MinLon != -180
}

func getLat2Coordinates(lat, degree float64) string {
	northSouth := 'S'
	if lat >= 0 {
		northSouth = 'N'
	}

	latPart := math.Abs(math.Round(lat/degree)) * degree

	if degree < 1 {
		return fmt.Sprintf("%s%02.1f", string(northSouth), latPart)
	}
	return fmt.Sprintf("%s%02.0f", string(northSouth), latPart)
}

func getLon2Coordinates(lon, degree float64) string {
	eastWest := 'W'
	if lon >= 0 {
		eastWest = 'E'
	}

	lonPart := math.Abs(math.Round(lon/degree)) * degree

	if degree < 1 {
		return fmt.Sprintf("%s%03.1f", string(eastWest), lonPart)
	}
	return fmt.Sprintf("%s%03.0f", string(eastWest), lonPart)
}

func CalculateTiles(coords gpx.BoundsType, degree float64) []string {
	var tiles []string

	lat1 := getLat2Coordinates(coords.MinLat, degree)
	lon1 := getLon2Coordinates(coords.MinLon, degree)
	lat2 := getLat2Coordinates(coords.MaxLat, degree)
	lon2 := getLon2Coordinates(coords.MaxLon, degree)

	tiles = append(tiles, lat1+lon1)
	if lon1 != lon2 {
		tiles = append(tiles, lat1+lon2)
	}
	if lat1 != lat2 {
		tiles = append(tiles, lat2+lon1)
	}
	if lat1 != lat2 && lon1 != lon2 {
		tiles = append(tiles, lat2+lon2)
	}

	return tiles
}

func GetPositionStart(g gpx.GPX) gpx.WptType {
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for _, WptType := range TrkSegType.TrkPt {
				if WptType.Lat != 0 && WptType.Lon != 0 {
					return *WptType
				}
			}
		}
	}
	return gpx.WptType{}
}

func GetPositionEnd(g gpx.GPX) gpx.WptType {
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for i := len(TrkSegType.TrkPt) - 1; i >= 0; i-- {
				WptType := TrkSegType.TrkPt[i]
				if WptType.Lat != 0 && WptType.Lon != 0 {
					return *WptType
				}
			}
		}
	}
	return gpx.WptType{}
}

func DistanceQuality(g gpx.GPX) float64 {
	var distance float64
	var quality float64 = 100
	var num int

	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for i := 0; i < len(TrkSegType.TrkPt)-1; i++ {
				distance += Distance2D(*TrkSegType.TrkPt[i], *TrkSegType.TrkPt[i+1])
			}
			num += len(TrkSegType.TrkPt)
		}
	}

	// check if points are very far away from each other
	step := distance / float64(num)
	if step > 30 {
		quality -= 12
	}
	if step > 8 {
		quality -= 6
	}

	// check intersections
	result := RemoveIntersections(g, 5, false)
	quality -= float64(len(result)) * 0.6

	// check first noise
	result = RemoveFirstNoise(g, false)
	quality -= float64(len(result)) * 0.3

	// check close points
	result = RemoveStops(g, 0.0, .5, math.MaxFloat64, 0, false)
	quality -= float64(len(result)) * 0.2

	// check high noise
	result = RemoveNoise(g, 6, 1.1, 4, false)
	quality -= float64(len(result)) * 0.4

	if quality < 0 {
		quality = 0
	}
	return quality
}
