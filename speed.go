package trackmaster

import (
	gpx "github.com/twpayne/go-gpx"
)

// MaxSpeed finds the max speed in a GPX file.
func MaxSpeed(g gpx.GPX, max float64, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if wptTypeNo != len(TrkSegType.TrkPt)-1 {
					point := SpeedBetween(*WptType, *TrkSegType.TrkPt[wptTypeNo+1], false)
					if point.Speed > max {
						maxSpeedFix(*TrkSegType, wptTypeNo, fix)
						point = SpeedBetween(*WptType, *TrkSegType.TrkPt[wptTypeNo+1], false)
						point.WptType = *TrkSegType.TrkPt[wptTypeNo]
						point.WptTypeNo = wptTypeNo
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

// SpeedBetween calculates the speed between two WptType.
func SpeedBetween(w, pt gpx.WptType, threeD bool) GPXElementInfo {
	seconds := TimeDiff(w, pt)
	var distLen float64
	if threeD {
		distLen = Distance3D(w, pt)
	} else {
		distLen = Distance2D(w, pt)
	}
	speed := distLen / seconds

	return GPXElementInfo{
		Speed:    speed,
		Length:   distLen,
		Duration: seconds,
	}
}

// maxSpeedFix fixes the max speed by adding a point in the middle of the two points.
func maxSpeedFix(ts gpx.TrkSegType, wptTypeNo int, fix bool) {
	if fix {
		closest := findClosestPoint(ts, wptTypeNo, 5)
		if closest == 0 {
			return
		}
		mid := midpoint(*ts.TrkPt[wptTypeNo], *ts.TrkPt[closest])
		ts.TrkPt[wptTypeNo+1].Lat = mid.Lat
		ts.TrkPt[wptTypeNo+1].Lon = mid.Lon
		ts.TrkPt[wptTypeNo+1].Ele = mid.Ele
	}
}

// findClosestPoint finds the closest point to a given point.
func findClosestPoint(ts gpx.TrkSegType, start, num int) int {
	var minDistance float64
	var minDistanceIndex int
	for i := start + 1; i < len(ts.TrkPt); i++ {
		num--
		if num == 0 {
			break
		}
		distance := Distance2D(*ts.TrkPt[start], *ts.TrkPt[i])
		if distance < minDistance || minDistance == 0 {
			minDistance = distance
			minDistanceIndex = i
		}
	}
	return minDistanceIndex
}
