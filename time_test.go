package trackmaster_test

import (
	"os"
	"testing"
	"time"

	trackmaster "github.com/inode64/gotrackmaster"
	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"
)

// testTimeFix tests the time fix function.
func TestTimeFix(t *testing.T) {
	filename := "testdata/carlos_prades_cool_de_la_creu.gpx"
	t.Run(filename, func(t *testing.T) {
		f, err := os.Open(filename)
		assert.NoError(t, err)
		defer f.Close()
		g, err := gpx.Read(f)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		n := trackmaster.FixTimesTrack(*g, true)
		assert.Equal(t, 39, n)
		datetest := time.Date(2015, time.April, 18, 7, 57, 51, 500000000, time.UTC)
		datetrack := g.Trk[0].TrkSeg[0].TrkPt[1].Time
		assert.Equal(t, datetrack, datetest)
		datetest = time.Date(2015, time.April, 18, 8, 4, 23, 0, time.UTC)
		datetrack = g.Trk[0].TrkSeg[0].TrkPt[4].Time
		assert.Equal(t, datetrack, datetest)
	})
}
