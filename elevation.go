package trackmaster

import (
	"math"

	gpx "github.com/twpayne/go-gpx"
)

// LostElevation finds the lost elevation
func LostElevation(g gpx.GPX, fix bool) []gpx.WptType {
	var result []gpx.WptType
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if wptTypeNo != len(TrkSegType.TrkPt)-1 {
					if WptType.Ele == 0 {
					    closest := findNextVerticalPoint(*TrkSegType, wptTypeNo, 10)
			    		if closest == -1 {
                            continue
                		}
					    TrkSegType.TrkPt[wptTypeNo].Ele = TrkSegType.TrkPt[closest].Ele
                        result = append(result, *TrkSegType.TrkPt[wptTypeNo])
					}
				}
			}
		}
	}
	return result
}


// MaxSpeedVertical finds the maximum vertical speed between two points.
func MaxSpeedVertical(g gpx.GPX, max float64, fix bool) []gpx.WptType {
	var result []gpx.WptType
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if wptTypeNo != len(TrkSegType.TrkPt)-1 {
					speed := SpeedVerticalBetween(*WptType, *TrkSegType.TrkPt[wptTypeNo+1])
					if speed > max {
						maxSpeedVerticalFix(*TrkSegType, wptTypeNo, fix)
						result = append(result, *TrkSegType.TrkPt[wptTypeNo])
					}
				}
			}
		}
	}
	return result
}

// SpeedVerticalBetween finds the vertical speed between two points.
func SpeedVerticalBetween(w, pt gpx.WptType) float64 {
	seconds := TimeDiff(w, pt)
	return ElevationAbs(w, pt) / seconds
}

// maxSpeedVerticalFix finds the maximum vertical speed between two points.
func maxSpeedVerticalFix(ts gpx.TrkSegType, wptTypeNo int, fix bool) {
	if fix {
		closest := findClosestVerticalPoint(ts, wptTypeNo, 5)
		if closest == 0 {
			return
		}
		mid := midpoint(*ts.TrkPt[wptTypeNo], *ts.TrkPt[closest])
		ts.TrkPt[wptTypeNo+1].Lat = mid.Lat
		ts.TrkPt[wptTypeNo+1].Lon = mid.Lon
		ts.TrkPt[wptTypeNo+1].Ele = mid.Ele
	}
}

// findClosestVerticalPoint finds the closest vertical point to the start point.
func findClosestVerticalPoint(ts gpx.TrkSegType, start, max int) int {
	var minElevation float64
	var minElevationIndex int
	var num int
	// find next closest point
	for i := start + 1; i < len(ts.TrkPt); i++ {
		num++
		if num > max {
			break
		}
		if ts.TrkPt[i].Ele == 0 {
            continue
        }
		elevation := ElevationAbs(*ts.TrkPt[start], *ts.TrkPt[i])
		if elevation < minElevation || minElevation == 0 {
			minElevation = elevation
			minElevationIndex = i
		}
	}

	return minElevationIndex
}

func findNextVerticalPoint(ts gpx.TrkSegType, start, max int) int {
	var num int
	// find next vertical point
	for i := start + 1; i < len(ts.TrkPt); i++ {
		num++
		if num > max {
			break
		}
		if ts.TrkPt[i].Ele != 0 {
            return i
        }
	}
	// find previous vertical point
	num = 0
	for i := start - 1; i > 0; i-- {
		num++
		if num > max {
			break
		}
		if ts.TrkPt[i].Ele != 0 {
            return i
        }
	}
	return -1
}

// Return the elevation of the midpoint between two points
func ElevationAbs(w, pt gpx.WptType) float64 {
	return w.Ele + math.Abs(w.Ele - pt.Ele) / 2
}
