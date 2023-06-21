package trackmaster

import (
	"math"

	"github.com/inode64/godem"
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

// RoundElevation rounds the elevation to 2 decimal places.
func RoundElevation(g gpx.GPX) {
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for _, WptType := range TrkSegType.TrkPt {
				WptType.Ele = math.Round(WptType.Ele*100) / 100
			}
		}
	}
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

func ElevationSRTM(g gpx.GPX) error {
	srtm, err := godem.NewSrtm(godem.SOURCE_ESA)
	if err != nil {
		return err
	}

	var hrs float64
	var lastHRS, lastLRS float64 // high ~ low resolution series

	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for wptTypeNo, WptType := range TrkSegType.TrkPt {
				elevation, _, err := srtm.GetElevation(WptType.Lat, WptType.Lon)
				if err != nil {
					return err
				}
				e := math.Abs(WptType.Ele - lastHRS)
				// fix only when the elevation is more than 10m different and the percentage is more than 3 meters
				// because the SRTM elevation is not very accurate, SRTM1 30 meters or SRTM3 90 meters
				if math.Abs(e) > 3 || lastLRS != elevation {
					hrs = 0
				}
				hrs += e

				lastHRS = WptType.Ele
				lastLRS = elevation

				TrkSegType.TrkPt[wptTypeNo].Ele = elevation
			}
		}
	}
	return nil
}

func ElevationSRTMAccuracy(g gpx.GPX) (int, error) {
	srtm, err := godem.NewSrtm(godem.SOURCE_ESA)
	if err != nil {
		return -1, err
	}

	var num, total int
	var max1, max2 float64

	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for _, WptType := range TrkSegType.TrkPt {
				elevation, _, err := srtm.GetElevation(WptType.Lat, WptType.Lon)
				if err != nil {
					return -1, err
				}
				max1 = 9
				max2 = 45
				if elevation > 250 {
					max1 = 8
					max2 = 40
				}
				if elevation > 500 {
					max1 = 6
					max2 = 35
				}
				if elevation > 1000 {
					max1 = 4
					max2 = 30
				}
				if elevation > 2000 {
					max1 = 3
					max2 = 20
				}
				if elevation > 3000 {
					max1 = 2
					max2 = 15
				}
				e := math.Abs(elevation-WptType.Ele) * 100 / elevation
				if e > max1 {
					num++
				}
				if e > max2 {
					num += 4
				}
				total++
			}
		}
	}
	if num > total {
		return 0, nil
	}
	if total == 0 {
		return 0, nil
	}
	return 100 - (num * 100 / total), nil
}
