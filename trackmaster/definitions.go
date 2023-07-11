package trackmaster

import (
	"errors"

	"github.com/sirupsen/logrus"
	gpx "github.com/twpayne/go-gpx"
)

const (
	earthRadius = 6371 * 1000
	oneDegree   = 1000.0 * 10000.8 / 90.0
)

type GPXElementInfo struct {
	TrkTypeNo     int
	WptTypeNo     int
	TrkSegTypeNo  int
	Count         int
	Length        float64
	Speed         float64
	SpeedVertical float64
	Elevation     float64
	Duration      float64
	WptType       gpx.WptType
}

const (
	ClassificationNone             = "Unknown"
	ClassificationCyClingSport     = "Cycling Sport"
	ClassificationCyClingMountain  = "Cycling Mountain"
	ClassificationCyClingTransport = "Cycling Transport"
	ClassificationCyClingTouring   = "Cycling Touring"
	ClassificationCyClingRacing    = "Cycling Racing"
	ClassificationCyClingIndoor    = "Cycling Indoor"
	ClassificationCyClingOther     = "Cycling Other"
	ClassificationRunningSport     = "Running Sport"
	ClassificationRunningMountain  = "Running Mountain"
	ClassificationRunningRacing    = "Running Racing"
	ClassificationRunningIndoor    = "Running Indoor"
	ClassificationRunningOther     = "Running Other"
	ClassificationWalkingSport     = "Walking Sport"
	ClassificationWalkingMountain  = "Walking Mountain"
	ClassificationWalkingTransport = "Walking Transport"
	ClassificationWalkingIndoor    = "Walking Indoor"
	ClassificationWalkingOther     = "Walking Other"
	ClassificationHikingSport      = "Hiking Sport"
	ClassificationHikingMountain   = "Hiking Mountain"
	ClassificationHikingOther      = "Hiking Other"
	ClassificationSwimmingSport    = "Swimming Sport"
	ClassificationSwimmingIndoor   = "Swimming Indoor"
	ClassificationRowingSport      = "Rowing Sport"
	ClassificationViaFerrataSport  = "Via Ferrata Sport"
	ClassificationMotorSport       = "Motor Sport"
)

const MinSegmentLength = 80

type Config struct {
	LogLevel logrus.Level
}

var Log = logrus.New()

var ErrNoLocation = errors.New("no location found")
