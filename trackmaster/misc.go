package trackmaster

import (
	"math"
	"os"

	"github.com/sirupsen/logrus"
	gpx "github.com/twpayne/go-gpx"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type MoveTrk struct {
	Track   int
	Segment int
}

func MoveSegment(g gpx.GPX, minPoints int, fix bool) []GPXElementInfo {
	var result []GPXElementInfo
	var move []MoveTrk
	dst := g
	for trkTypeNo, TrkType := range g.Trk {
		if len(TrkType.TrkSeg) < 2 {
			continue
		}
		for trkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			if len(TrkSegType.TrkPt) > minPoints {
				continue
			}
			move = append(move, MoveTrk{trkTypeNo, trkSegTypeNo})
			if len(TrkSegType.TrkPt) == 0 {
				continue
			}
			pre := CompareTime(g, trkTypeNo, trkSegTypeNo, false)
			next := CompareTime(g, trkTypeNo, trkSegTypeNo, true)
			var point GPXElementInfo
			if pre < next {
				point = GPXElementInfo{
					WptType:      *TrkSegType.TrkPt[0],
					WptTypeNo:    0,
					TrkSegTypeNo: trkSegTypeNo,
					TrkTypeNo:    trkTypeNo,
					Count:        len(TrkSegType.TrkPt),
				}
				dst.Trk[trkTypeNo].TrkSeg[trkSegTypeNo-1].TrkPt = append(dst.Trk[trkTypeNo].TrkSeg[trkSegTypeNo-1].TrkPt, TrkSegType.TrkPt...)
			} else {
				point = GPXElementInfo{
					WptType:      *TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1],
					WptTypeNo:    len(TrkSegType.TrkPt) - 1,
					TrkSegTypeNo: trkSegTypeNo,
					TrkTypeNo:    trkTypeNo,
					Count:        len(TrkSegType.TrkPt),
				}
				dst.Trk[trkTypeNo].TrkSeg[trkSegTypeNo+1].TrkPt = append(TrkSegType.TrkPt, dst.Trk[trkTypeNo].TrkSeg[trkSegTypeNo+1].TrkPt...)
			}
			result = append(result, point)
		}
	}
	if fix {
		for i := len(move) - 1; i >= 0; i-- {
			if move[i].Segment == 0 {
				dst.Trk[move[i].Track].TrkSeg = dst.Trk[move[i].Track].TrkSeg[1:]
			} else if move[i].Segment == len(dst.Trk[move[i].Track].TrkSeg)-1 {
				dst.Trk[move[i].Track].TrkSeg = dst.Trk[move[i].Track].TrkSeg[:len(dst.Trk[move[i].Track].TrkSeg)-1]
			} else {
				dst.Trk[move[i].Track].TrkSeg = append(dst.Trk[move[i].Track].TrkSeg[:move[i].Segment], dst.Trk[move[i].Track].TrkSeg[move[i].Segment+1:]...)
			}
		}
		g = dst
	}
	return result
}

func CompareTime(g gpx.GPX, trkTypeNo, trkSegTypeNo int, end bool) float64 {
	if end {
		p := *g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt[len(g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt)-1]
		trkTypeNo, trkSegTypeNo = NextSegment(g, trkTypeNo, trkSegTypeNo)
		if trkTypeNo == -1 {
			return math.MaxFloat64
		}
		return TimeDiff(p, *g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt[0])
	}
	p := *g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt[0]
	trkTypeNo, trkSegTypeNo = PreviousSegment(g, trkTypeNo, trkSegTypeNo)
	if trkTypeNo == -1 {
		return math.MaxFloat64
	}
	return TimeDiff(p, *g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt[len(g.Trk[trkTypeNo].TrkSeg[trkSegTypeNo].TrkPt)-1])
}

func NextSegment(g gpx.GPX, trkTypeNo, trkSegTypeNo int) (int, int) {
	if trkSegTypeNo >= len(g.Trk[trkTypeNo].TrkSeg)-1 {
		trkSegTypeNo = 0
		trkTypeNo++
		if trkTypeNo >= len(g.Trk)-1 {
			return -1, -1
		}
	} else {
		trkSegTypeNo++
	}
	return trkTypeNo, trkSegTypeNo
}

func PreviousSegment(g gpx.GPX, trkTypeNo, trkSegTypeNo int) (int, int) {
	if trkSegTypeNo <= 0 {
		trkSegTypeNo = len(g.Trk[trkTypeNo].TrkSeg) - 1
		trkTypeNo--
		if trkTypeNo <= 0 {
			return -1, -1
		}
	} else {
		trkSegTypeNo--
	}
	return trkTypeNo, trkSegTypeNo
}

func ClassificationTrack(filename string) string {
	var speedUp, speedDown, speedFlat, speedTotal, elevation, distance float64
	var total int

	f, err := os.Open(filename)
	if err != nil {
		return ClassificationNone
	}
	defer f.Close()

	g, err := gpx.Read(f)
	if err != nil {
		return ClassificationNone
	}

	// first the points without time are corrected, because it is needed for the rest of the functions
	_ = FixTimesTrack(*g, true)

	// Removes points that have been recorded excessively far away, at more than 200 m/s
	_ = MaxSpeed(*g, 200, true)

	// Simplifies the track and removes points that are not necessary
	_ = RemoveStops(*g, 0.0, 1.2, math.MaxFloat64, 0, true)

	// We remove the stops of more than 90 seconds in less than 5 meters
	_ = RemoveStops(*g, 30.0, 9.0, 8, 12, true)

	CheckIntersecting(*g, 7, true)
	CheckIntersecting(*g, 7, true)
	CheckIntersecting(*g, 7, true)
	CheckIntersecting(*g, 7, true)

	num, err := ElevationSRTMAccuracy(*g)
	if err != nil {
		if num < 60 {
			_ = ElevationSRTM(*g)
		}
	}

	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			if len(TrkSegType.TrkPt) < MinSegmentLength {
				continue
			}
			div := len(TrkSegType.TrkPt) / 10
			// only check middle 80.0% of track
			for i := div; i < len(TrkSegType.TrkPt)-div; i++ {
				point := SpeedBetween(*TrkSegType.TrkPt[i], *TrkSegType.TrkPt[i+1], false)
				if point.SpeedVertical <= 0.4 {
					speedFlat += point.Speed
				}
				if point.SpeedVertical > 0.4 {
					speedUp += point.Speed
				}
				if point.SpeedVertical < -0.4 {
					speedDown += point.Speed
				}
				speedTotal += point.Speed
				elevation += math.Abs(point.Elevation)
				distance += point.Length

				total++
			}
		}
	}

	speedUp /= float64(total)
	speedDown /= float64(total)
	speedFlat /= float64(total)
	speedTotal /= float64(total)

	c := ClassificationNone

	if total != 0 {
		// Flat sports
		if (elevation / distance) < 0.05 {
			c = ClassificationWalkingTransport
			if speedFlat > 1.6 {
				c = ClassificationRunningSport
			}
			if speedFlat > 4.1 {
				c = ClassificationCyClingTransport
			}
			if speedFlat > 7.5 {
				c = ClassificationCyClingSport
			}
			if speedFlat > 11 {
				c = ClassificationCyClingRacing
			}
			if speedFlat > 25 {
				c = ClassificationMotorSport
			}
		} else {
			c = ClassificationWalkingMountain
			/* TODO: Check this
			if speedDown < 0.1 && speedTotal < 0.5 {
				c = ClassificationViaFerrataSport
			}
			*/
			if speedFlat > 1.2 || speedTotal > 1.3 {
				c = ClassificationRunningMountain
			}
			if speedFlat > 3.8 || speedTotal > 3.8 {
				c = ClassificationCyClingMountain
			}
		}
	}

	Log.WithFields(logrus.Fields{
		"Elevation":                          elevation,
		"Ratio of elevation versus distance": elevation / distance,
		"Upload speed":                       speedUp,
		"Lowering speed":                     speedDown,
		"Flat speed":                         speedFlat,
		"Average speed":                      speedTotal,
		"Classification":                     c,
	}).Debug("Classification result")

	return c
}
