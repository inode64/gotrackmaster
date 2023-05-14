package trackmaster

import (
	"math"

	gpx "github.com/twpayne/go-gpx"
)

// LostElevation finds the lost elevation.
func LostElevation(g gpx.GPX, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if WptType.Ele == 0 {
					closest := findNextVerticalPoint(*TrkSegType, wptTypeNo, 10)
					if closest == -1 {
						continue
					}
					if fix {
						TrkSegType.TrkPt[wptTypeNo].Ele = TrkSegType.TrkPt[closest].Ele
					}
					point := GPXElementInfo{}
					point.WptType = *TrkSegType.TrkPt[wptTypeNo]
					point.wptTypeNo = wptTypeNo
					point.TrkSegTypeNo = TrkSegTypeNo
					point.TrkTypeNo = TrkTypeNo
					point.Elevation = TrkSegType.TrkPt[closest].Ele

					result = append(result, point)
				}
			}
		}
	}
	return result
}

// MaxSpeedVertical finds the maximum vertical speed between two points.
func MaxSpeedVertical(g gpx.GPX, max float64, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if wptTypeNo != len(TrkSegType.TrkPt)-1 {
					point := SpeedVerticalBetween(*WptType, *TrkSegType.TrkPt[wptTypeNo+1])
					if point.Speed > max {
						maxSpeedVerticalFix(*TrkSegType, wptTypeNo, fix)
						point.WptType = *TrkSegType.TrkPt[wptTypeNo]
						point.wptTypeNo = wptTypeNo
						point.TrkSegTypeNo = TrkSegTypeNo
						point.TrkTypeNo = TrkTypeNo

						result = append(result, point)
					}
				}
			}
		}
	}
	return result
}

// SpeedVerticalBetween finds the vertical speed between two points.
func SpeedVerticalBetween(w, pt gpx.WptType) GPXElementInfo {
	seconds := TimeDiff(w, pt)
	elevation := ElevationAbs(w, pt)
	speed := elevation / seconds

	return GPXElementInfo{
		Speed:    speed,
		Length:   elevation,
		Duration: seconds,
	}
}

// maxSpeedVerticalFix finds the maximum vertical speed between two points.
func maxSpeedVerticalFix(ts gpx.TrkSegType, wptTypeNo int, fix bool) {
	if fix {
		closest := findClosestVerticalPoint(ts, wptTypeNo, 5)
		if closest == 0 {
			return
		}
		ts.TrkPt[wptTypeNo+1].Ele = MiddleElevation(*ts.TrkPt[wptTypeNo], *ts.TrkPt[closest])
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
		elevation := MiddleElevation(*ts.TrkPt[start], *ts.TrkPt[i])
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

// Return the elevation of the midpoint between two points.
func ElevationAbs(w, pt gpx.WptType) float64 {
	return math.Abs(w.Ele - pt.Ele)
}

func MiddleElevation(w, pt gpx.WptType) float64 {
	return pt.Ele + (w.Ele-pt.Ele)/2
}
