package trackmaster_test

import (
	"os"
	"testing"

	trackmaster "github.com/inode64/gotrackmaster"
	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"
)

// testSpeedFix tests the speed fix algorithm.
func TestLostElevation(t *testing.T) {
	filename := "testdata/2020-12-19_11-14_Sat_benitandus.gpx"
	t.Run(filename, func(t *testing.T) {
		f, err := os.Open(filename)
		assert.NoError(t, err)
		defer f.Close()
		g, err := gpx.Read(f)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		trackmaster.LostElevation(*g, true)
		assert.Equal(t, g.Trk[0].TrkSeg[0].TrkPt[3655].Ele, 468.95)
	})
}
