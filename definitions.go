package trackmaster

import (
	gpx "github.com/twpayne/go-gpx"
)

const (
	earthRadius = 6371 * 1000
	oneDegree   = 1000.0 * 10000.8 / 90.0
)

type GPXElementInfo struct {
	TrkTypeNo    int
	WptTypeNo    int
	TrkSegTypeNo int
	Count        int
	Length       float64
	Speed        float64
	Elevation    float64
	Duration     float64
	WptType      gpx.WptType
}
