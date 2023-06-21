package trackmaster_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"

	trackmaster "github.com/inode64/gotrackmaster/trackmaster"
)

// testSpeedFix tests the speed fix algorithm.
func TestSpeedFix(t *testing.T) {
	filename := "testdata/2020-10-03_09-05_Sat_pedraforca.gpx"
	t.Run(filename, func(t *testing.T) {
		f, err := os.Open(filename)
		assert.NoError(t, err)
		defer f.Close()
		g, err := gpx.Read(f)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		trackmaster.MaxSpeed(*g, 300, true)
		w := gpx.WptType{Lat: 42.24870745000008, Lon: 1.664240950000083}
		assert.Equal(t, g.Trk[0].TrkSeg[0].TrkPt[504].Lat, w.Lat)
		assert.Equal(t, g.Trk[0].TrkSeg[0].TrkPt[504].Lon, w.Lon)
		w = gpx.WptType{Lat: 42.2516829000064, Lon: 1.6696103500038912}
		assert.Equal(t, g.Trk[0].TrkSeg[0].TrkPt[662].Lat, w.Lat)
		assert.Equal(t, g.Trk[0].TrkSeg[0].TrkPt[662].Lon, w.Lon)
		w = gpx.WptType{Lat: 42.24071075000116, Lon: 1.7195158000128756}
		assert.Equal(t, g.Trk[1].TrkSeg[0].TrkPt[7].Lat, w.Lat)
		assert.Equal(t, g.Trk[1].TrkSeg[0].TrkPt[7].Lon, w.Lon)
	})
}
