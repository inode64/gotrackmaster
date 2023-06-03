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
func FixTimesSegment(tr gpx.TrkSegType) (gpx.TrkSegType, int) {
	var num int
	var lastValidTime time.Time
	if len(tr.TrkPt) < 2 {
		return tr, num
	}
	// Find the first valid time
	for _, WptType := range tr.TrkPt {
		if lastValidTime.IsZero() || WptType.Time.Before(lastValidTime) {
			lastValidTime = WptType.Time
		}
	}
	// Check first element
	if !tr.TrkPt[0].Time.IsZero() && tr.TrkPt[0].Time.After(tr.TrkPt[1].Time) {
		tr.TrkPt[0].Time = tr.TrkPt[1].Time.Add(-10 * time.Second)
		num++
	}
	// Check all intermediate elements
	lastValidTime = tr.TrkPt[0].Time
	for i := 1; i < len(tr.TrkPt)-1; i++ {
		if tr.TrkPt[i].Time.IsZero() {
			continue
		}
		maxValidTime := lastValidTime.Add(time.Hour)
		if tr.TrkPt[i].Time.After(tr.TrkPt[i+1].Time) || tr.TrkPt[i].Time.After(maxValidTime) {
			tr.TrkPt[i].Time = findNextValidTime(tr, lastValidTime, i)
			num++
		} else {
			lastValidTime = tr.TrkPt[i].Time
		}
	}
	return tr, num
}

func findNextValidTime(tr gpx.TrkSegType, lastValidTime time.Time, start int) time.Time {
	maxValidTime := lastValidTime.Add(time.Hour)

	for i := start + 1; i < len(tr.TrkPt); i++ {
		if tr.TrkPt[i].Time.IsZero() {
			continue
		}
		if tr.TrkPt[i].Time.After(lastValidTime) && tr.TrkPt[i].Time.Before(maxValidTime) {
			return lastValidTime.Add(tr.TrkPt[i].Time.Sub(lastValidTime) / time.Duration(i-start+1))
		}
	}

	// check when there is no valid time
	return tr.TrkPt[0].Time
}

// FixTimesTrack fixes the time of a track.
func FixTimesTrack(g gpx.GPX, fix bool) int {
	var num, n int
	for TrkTypeNo, TrkType := range g.Trk {
		for TrkSegTypeNo, TrkSegType := range TrkType.TrkSeg {
			if fix {
				*g.Trk[TrkTypeNo].TrkSeg[TrkSegTypeNo], n = FixTimesSegment(*TrkSegType)
			} else {
				_, n = FixTimesSegment(*TrkSegType)
			}
			num += n
		}
	}
	return num
}

// TimeEmpty returns true if there is no time information in the GPX file.
func TimeEmpty(g gpx.GPX) bool {
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			for _, WptType := range TrkSegType.TrkPt {
				if !WptType.Time.IsZero() {
					return false
				}
			}

		}
	}
	return true
}

// TimeQuality returns the quality of the time information in the GPX file.
func TimeQuality(g gpx.GPX) int {
	var num, total int
	for _, TrkType := range g.Trk {
		for _, TrkSegType := range TrkType.TrkSeg {
			var lastValidTime time.Time
			for _, WptType := range TrkSegType.TrkPt {
				if !WptType.Time.IsZero() {
					num++
				}
				if !lastValidTime.IsZero() && WptType.Time.Before(lastValidTime) {
					num += 4
				}
				lastValidTime = WptType.Time
				total++
			}
		}
	}
	if num > total {
		return 0
	}
	if total == 0 {
		return -1
	}
	return 100 - (num * 100 / total)
}
