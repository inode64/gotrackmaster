package trackmaster

import (
	"math"
	"net/http"

	"github.com/tkrajina/go-elevations/geoelevations"
	gpx "github.com/twpayne/go-gpx"
)

// LostElevation finds the lost elevation.
func LostElevation(g gpx.GPX, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if WptType.Ele <= 0 {
					closest := findNextVerticalPoint(*TrkSegType, wptTypeNo, 10)
					if closest == -1 {
						continue
					}
					point := GPXElementInfo{}
					point.WptType = *TrkSegType.TrkPt[wptTypeNo]
					point.WptTypeNo = wptTypeNo
					point.TrkSegTypeNo = TrkSegTypeNo
					point.TrkTypeNo = TrkTypeNo
					point.Elevation = TrkSegType.TrkPt[closest].Ele

					result = append(result, point)

					if fix {
						TrkSegType.TrkPt[wptTypeNo].Ele = TrkSegType.TrkPt[closest].Ele
					}
				}
			}
		}
	}
	return result
}

func ContinuousElevation(g gpx.GPX, count int, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	var lastElevation float64
	var num, start, end int
	var point GPXElementInfo

	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				if lastElevation != WptType.Ele {
					if num > count {
						point.Count = start - end + 1
						result = append(result, point)
						if fix {
							continuousElevationFix(*TrkSegType, start, end)
						}
					}
					end = 0
					num = 0
					start = wptTypeNo
				}
				if num >= count {
					if end == 0 {
						point = GPXElementInfo{}
						point.WptType = *TrkSegType.TrkPt[wptTypeNo]
						point.WptTypeNo = wptTypeNo
						point.TrkSegTypeNo = TrkSegTypeNo
						point.TrkTypeNo = TrkTypeNo
					}
					end = wptTypeNo
				}
				num++
				lastElevation = WptType.Ele
				end = wptTypeNo
			}
			if num > count {
				point.Count = start - end + 1
				result = append(result, point)
				if fix {
					continuousElevationFix(*TrkSegType, start, end)
				}
			}

		}
	}
	return result
}

func continuousElevationFix(ts gpx.TrkSegType, start, end int) {
	srtm, err := geoelevations.NewSrtm(http.DefaultClient)
	if err != nil {
		return
	}

	for i := start; i < end; i++ {
		elevation, err := srtm.GetElevation(http.DefaultClient, ts.TrkPt[i].Lat, ts.TrkPt[i].Lon)
		if err != nil || elevation == 0 {
			continue
		}
		ts.TrkPt[i].Ele = elevation
	}
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
						point.WptType = *TrkSegType.TrkPt[wptTypeNo]
						point.WptTypeNo = wptTypeNo
						point.TrkSegTypeNo = TrkSegTypeNo
						point.TrkTypeNo = TrkTypeNo
						result = append(result, point)

						if fix {
							gaussianFilter(*TrkSegType, wptTypeNo-2, wptTypeNo+5, 3, 1.5)
						}
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

/*
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
}*/

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

func ElevationSRTM(g gpx.GPX, max float64, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	srtm, err := geoelevations.NewSrtm(http.DefaultClient)
	if err != nil {
		return result
	}

	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				elevation, err := srtm.GetElevation(http.DefaultClient, WptType.Lat, WptType.Lon)
				if err != nil || elevation == 0 {
					continue
				}
				e := math.Abs(WptType.Ele - elevation)
				p := e * 100 / WptType.Ele
				// fix only when the elevation is more than 10m different and the percentage is more than max
				// because the STRM elevation is not very accurate, STRM1 30 meters, STRM3 90 meters
				if p > max && e > 10 {
					point := GPXElementInfo{}
					point.WptType = *TrkSegType.TrkPt[wptTypeNo]
					point.WptTypeNo = wptTypeNo
					point.TrkSegTypeNo = TrkSegTypeNo
					point.TrkTypeNo = TrkTypeNo
					point.Elevation = elevation

					result = append(result, point)

					if fix {
						TrkSegType.TrkPt[wptTypeNo].Ele = elevation
					}
				}
			}
		}
	}
	return result
}
