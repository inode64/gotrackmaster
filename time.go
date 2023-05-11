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
