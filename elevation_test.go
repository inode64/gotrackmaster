package trackmaster_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"

	trackmaster "github.com/inode64/gotrackmaster"
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

func TestLostElevationSRTM(t *testing.T) {
	filename := "testdata/2023-03-05_09-27_Sun.gpx"
	t.Run(filename, func(t *testing.T) {
		f, err := os.Open(filename)
		assert.NoError(t, err)
		defer f.Close()
		g, err := gpx.Read(f)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		err = trackmaster.ElevationSRTM(*g)
		assert.NoError(t, err)
		assert.Equal(t, g.Trk[0].TrkSeg[2].TrkPt[265].Ele, 721.0)
		assert.Equal(t, g.Trk[0].TrkSeg[2].TrkPt[601].Ele, 852.0)
	})
}
