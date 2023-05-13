package trackmaster

import (
	"time"

	gpx "github.com/twpayne/go-gpx"
)

// TimeDiff returns the time difference of two WptType in seconds.
func TimeDiff(w, pt gpx.WptType) float64 {
	t1 := w.Time
	t2 := pt.Time
	if t1.Equal(t2) {
		return 0.0
	}
	var delta time.Duration
	if t1.After(t2) {
		delta = t1.Sub(t2)
	} else {
		delta = t2.Sub(t1)
	}
	return delta.Seconds()
}

// FixTimesSegment fixes the time of a track segment.
func FixTimesSegment(tr gpx.TrkSegType) gpx.TrkSegType {
	if len(tr.TrkPt) == 0 {
		return tr
	}
	// Check first element
	if !tr.TrkPt[0].Time.IsZero() && tr.TrkPt[0].Time.After(tr.TrkPt[1].Time) {
		tr.TrkPt[0].Time = tr.TrkPt[1].Time.Add(-10 * time.Second)
	}
	// Check all intermediate elements
	lastValidTime := tr.TrkPt[0].Time
	totalDiff := time.Duration(0)
	for i := 1; i < len(tr.TrkPt)-1; i++ {
		if tr.TrkPt[i].Time.IsZero() {
			continue
		}
		if tr.TrkPt[i].Time.After(tr.TrkPt[i+1].Time) || tr.TrkPt[i].Time.Before(lastValidTime) {
			averageDiff := totalDiff / time.Duration(i)
			tr.TrkPt[i].Time = lastValidTime.Add(averageDiff)
		} else {
			totalDiff += tr.TrkPt[i].Time.Sub(lastValidTime)
			lastValidTime = tr.TrkPt[i].Time
		}
	}
	// Check last element
	if tr.TrkPt[len(tr.TrkPt)-1].Time.Before(lastValidTime) {
		tr.TrkPt[len(tr.TrkPt)-1].Time = lastValidTime.Add(totalDiff / time.Duration(len(tr.TrkPt)-1))
	}
	return tr
}

// FixTimesTrack fixes the time of a track.
func FixTimesTrack(g gpx.GPX) {
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			*g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo] = FixTimesSegment(*TrkSegType)
		}
	}
}
