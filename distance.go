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
	var lastPoint int = -1
	var minDistance float64 = math.MaxFloat64
	for i := start + 1; i < start+max; i++ {
		if i >= len(ts.TrkPt) {
			break
		}
		distance := HaversineDistanceTrkPt(*ts.TrkPt[start], *ts.TrkPt[i])
		elevation := ElevationAbs(*ts.TrkPt[start], *ts.TrkPt[i])

		if distance < minDistance && distance < maxDistance && elevation <= maxElevation {
			minDistance = distance
			lastPoint = i
		}
	}

	if lastPoint == -1 {
		return -1, minDistance
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

// remove points when accuracy is too low in first point
func RemoveFirstNoise(g gpx.GPX, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			var dst []*gpx.WptType
			for i := 0; i < 11; i++ {
				if i+1 >= len(TrkSegType.TrkPt) {
					if fix {
						dst = append(dst, TrkSegType.TrkPt[i])
					}
					break
				}

				nextDistance := HaversineDistanceTrkPt(*TrkSegType.TrkPt[i], *TrkSegType.TrkPt[i+1])
				closerPoint, closerDistance := findNextCloserPoint(*TrkSegType, i, 5, 8, 0)
				if nextDistance > closerDistance {
					point := GPXElementInfo{}
					point.WptType = *TrkSegType.TrkPt[i]
					point.WptTypeNo = i
					point.TrkSegTypeNo = TrkSegTypeNo
					point.TrkTypeNo = TrkTypeNo
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
				} else {
					if fix {
						if i >= 10 {
							dst = append(dst, TrkSegType.TrkPt[i:]...)
						} else {
							dst = append(dst, TrkSegType.TrkPt[i])
						}
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

func RemoveStops(g gpx.GPX, minSeconds, maxDistance float64, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	var distance float64
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			var dst []*gpx.WptType
			var firstPoint int = -1
			var numPoints, point int
			for wptTypeNo, _ := range TrkSegType.TrkPt {
				if wptTypeNo == len(TrkSegType.TrkPt)-1 {
					continue
				}

				if firstPoint == -1 {
					point = wptTypeNo
				} else {
					point = firstPoint
				}

				distance = HaversineDistanceTrkPt(*TrkSegType.TrkPt[point], *TrkSegType.TrkPt[wptTypeNo+1])
				if distance <= maxDistance {
					if firstPoint == -1 {
						firstPoint = wptTypeNo
					}
					numPoints++
				} else {
					seconds := TimeDiff(*TrkSegType.TrkPt[point], *TrkSegType.TrkPt[wptTypeNo])
					if fix {
						if numPoints > 3 && seconds > minSeconds {
							dst = append(dst, TrkSegType.TrkPt[firstPoint])
						} else {
							if firstPoint != -1 {
								for i := firstPoint; i < wptTypeNo; i++ {
									dst = append(dst, TrkSegType.TrkPt[i])
								}
							}
						}
						dst = append(dst, TrkSegType.TrkPt[wptTypeNo])
					}
					if firstPoint != -1 && numPoints > 3 && seconds > minSeconds {
						point := GPXElementInfo{}
						point.WptType = *TrkSegType.TrkPt[firstPoint]
						point.WptTypeNo = firstPoint
						point.TrkSegTypeNo = TrkSegTypeNo
						point.TrkTypeNo = TrkTypeNo
						result = append(result, point)
					}
					firstPoint = -1
					numPoints = 0
				}
			}
			if fix {
				if numPoints > 3 {
					dst = append(dst, TrkSegType.TrkPt[firstPoint])
				} else {
					for i := len(TrkSegType.TrkPt) - numPoints - 1; i < len(TrkSegType.TrkPt); i++ {
						dst = append(dst, TrkSegType.TrkPt[i])
					}
				}
				if firstPoint != -1 && numPoints > 3 {
					point := GPXElementInfo{}
					point.WptType = *TrkSegType.TrkPt[firstPoint]
					point.WptTypeNo = firstPoint
					point.TrkSegTypeNo = TrkSegTypeNo
					point.TrkTypeNo = TrkTypeNo
					result = append(result, point)
				}
				g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt = dst
			}
		}
	}

	return result
}

/*

// Function to smooth the GPS data using Gaussian filter
// https://github.com/gonum/gonum
// "gonum.org/v1/gonum/stat/distuv"

type point struct {
	lat float64
	lng float64
}

func smoothGPSData(data []point, sigma float64) []point {
	smoothedData := make([]point, len(data))
	x := make([]float64, len(data))
	y := make([]float64, len(data))

	// Extract the latitude and longitude data into separate slices
	for i, p := range data {
		x[i], y[i] = p.lat, p.lng
	}

	// Create a Gaussian distribution with sigma as the standard deviation
	dist := distuv.Normal{
		Mu:    0,
		Sigma: sigma,
	}

	// Convolve the data with the Gaussian kernel
	for i := range x {
		for j := range y {
			kernel := dist.Prob(x[i] - x[j])
			smoothedData[i].lat += kernel * x[j]
			smoothedData[i].lng += kernel * y[j]
		}
	}

	return smoothedData
}

// a machine learning algorithm to smooth horizontally
// "github.com/salkj/kmeans"
func smoothGPSPositions(coordinates [][]float64, numClusters int) [][]float64 {
	kmeans := kmeans.New()
	clusters, _ := kmeans.Clusterize(coordinates, numClusters)
	smoothedPositions := make([][]float64, numClusters)
	for i, cluster := range clusters {
		avgLat, avgLon := 0.0, 0.0
		for _, point := range cluster.Points {
			avgLat += point[0]
			avgLon += point[1]
		}
		avgLat /= float64(len(cluster.Points))
		avgLon /= float64(len(cluster.Points))
		smoothedPositions[i] = []float64{avgLat, avgLon}
	}
	return smoothedPositions
}

// https://jeffreyearly.com/smoothing-and-interpolating-noisy-gps-data/
// b-spline interpolation
func smoothGPSData2(data []point, sigma float64) []point {
	smoothedData := make([]point, len(data))
	x := make([]float64, len(data))
	y := make([]float64, len(data))

	// Extract the latitude and longitude data into separate slices
	for i, p := range data {
		x[i], y[i] = p.lat, p.lng
	}

	// Create a Gaussian distribution with sigma as the standard deviation
	dist := distuv.Normal{
		Mu:    0,
		Sigma: sigma,
	}

	// Convolve the data with the Gaussian kernel
	for i := range x {
		for j := range y {
			kernel := dist.Prob(x[i] - x[j])
			smoothedData[i].lat += kernel * x[j]
			smoothedData[i].lng += kernel * y[j]
		}
	}

	return smoothedData
}
*/
