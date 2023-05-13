package trackmaster_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"

	trackmaster "github.com/inode64/gotrackmaster"
)

// testTimeFix tests the time fix function.
func TestTimeFix(t *testing.T) {
	t.Run("testdata/carlos_prades_cool_de_la_creu.gpx", func(t *testing.T) {
		f, err := os.Open("testdata/carlos_prades_cool_de_la_creu.gpx")
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
