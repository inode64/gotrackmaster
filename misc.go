package trackmaster

import (
	"math"

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
	for TrkTypeNo, TrkType := range g.Trk {
		if len(TrkType.TrkSeg) < 2 {
			continue
		}
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			if len(TrkSegType.TrkPt) > minPoints {
				continue
			}
			move = append(move, MoveTrk{TrkTypeNo, TrkSegTypeNo})
			if len(TrkSegType.TrkPt) == 0 {
				continue
			}
			pre := CompareTime(g, TrkTypeNo, TrkSegTypeNo, false)
			next := CompareTime(g, TrkTypeNo, TrkSegTypeNo, true)
			var point GPXElementInfo
			if pre < next {
				point = GPXElementInfo{
					WptType:      *TrkSegType.TrkPt[0],
					WptTypeNo:    0,
					TrkSegTypeNo: TrkSegTypeNo,
					TrkTypeNo:    TrkTypeNo,
					Count:        len(TrkSegType.TrkPt),
				}
				dst.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo-1].TrkPt = append(dst.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo-1].TrkPt, TrkSegType.TrkPt...)
			} else {
				point = GPXElementInfo{
					WptType:      *TrkSegType.TrkPt[len(TrkSegType.TrkPt)-1],
					WptTypeNo:    len(TrkSegType.TrkPt) - 1,
					TrkSegTypeNo: TrkSegTypeNo,
					TrkTypeNo:    TrkTypeNo,
					Count:        len(TrkSegType.TrkPt),
				}
				dst.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo+1].TrkPt = append(TrkSegType.TrkPt, dst.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo+1].TrkPt...)
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

func CompareTime(g gpx.GPX, TrkTypeNo, TrkSegTypeNo int, end bool) float64 {
	if end {
		p := *g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt[len(g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt)-1]
		TrkTypeNo, TrkSegTypeNo = NextSegment(g, TrkTypeNo, TrkSegTypeNo)
		if TrkTypeNo == -1 {
			return math.MaxFloat64
		}
		return TimeDiff(p, *g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt[0])
	}
	p := *g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt[0]
	TrkTypeNo, TrkSegTypeNo = PreviousSegment(g, TrkTypeNo, TrkSegTypeNo)
	if TrkTypeNo == -1 {
		return math.MaxFloat64
	}
	return TimeDiff(p, *g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt[len(g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo].TrkPt)-1])
}

func NextSegment(g gpx.GPX, TrkTypeNo, TrkSegTypeNo int) (int, int) {
	if TrkSegTypeNo >= len(g.Trk[TrkTypeNo].TrkSeg)-1 {
		TrkSegTypeNo = 0
		TrkTypeNo++
		if TrkTypeNo >= len(g.Trk)-1 {
			return -1, -1
		}
	} else {
		TrkSegTypeNo++
	}
	return TrkTypeNo, TrkSegTypeNo
}

func PreviousSegment(g gpx.GPX, TrkTypeNo, TrkSegTypeNo int) (int, int) {
	if TrkSegTypeNo <= 0 {
		TrkSegTypeNo = len(g.Trk[TrkTypeNo].TrkSeg) - 1
		TrkTypeNo--
		if TrkTypeNo <= 0 {
			return -1, -1
		}
	} else {
		TrkSegTypeNo--
	}
	return TrkTypeNo, TrkSegTypeNo
}
