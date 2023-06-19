package trackmaster

import (
	gpx "github.com/twpayne/go-gpx"

	"github.com/sirupsen/logrus"
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

const ClassificationNone = "Unknown"
const ClassificationCyClingSport = "Cycling Sport"
const ClassificationCyClingMountain = "Cycling Mountain"
const ClassificationCyClingTransport = "Cycling Transport"
const ClassificationCyClingTouring = "Cycling Touring"
const ClassificationCyClingRacing = "Cycling Racing"
const ClassificationCyClingIndoor = "Cycling Indoor"
const ClassificationCyClingOther = "Cycling Other"
const ClassificationRunningSport = "Running Sport"
const ClassificationRunningMountain = "Running Mountain"
const ClassificationRunningRacing = "Running Racing"
const ClassificationRunningIndoor = "Running Indoor"
const ClassificationRunningOther = "Running Other"
const ClassificationWalkingSport = "Walking Sport"
const ClassificationWalkingMountain = "Walking Mountain"
const ClassificationWalkingTransport = "Walking Transport"
const ClassificationWalkingIndoor = "Walking Indoor"
const ClassificationWalkingOther = "Walking Other"
const ClassificationHikingSport = "Hiking Sport"
const ClassificationHikingMountain = "Hiking Mountain"
const ClassificationHikingOther = "Hiking Other"
const ClassificationSwimmingSport = "Swimming Sport"
const ClassificationSwimmingIndoor = "Swimming Indoor"
const ClassificationRowingSport = "Rowing Sport"
const ClassificationViaFerrataSport = "Via Ferrata Sport"
const ClassificationMotorSport = "Motor Sport"

const MinSegmentLength = 80

type Config struct {
	LogLevel logrus.Level
}

var Log = logrus.New()

/*
func Init(config *Config) {
	if config == nil {
		config.LogLevel = logrus.WarnLevel
	}

	Log.SetLevel(config.LogLevel)
}
*/
