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

const ClassificationNone = "Unkonwn"
const ClassificatiomCyClingSport = "Cycling Sport"
const ClassificatiomCyClingMountain = "Cycling Mountain"
const ClassificatiomCyClingTransport = "Cycling Transport"
const ClassificatiomCyClingTouring = "Cycling Touring"
const ClassificatiomCyClingRacing = "Cycling Racing"
const ClassificatiomCyClingIndoor = "Cycling Indoor"
const ClassificatiomCyClingOther = "Cycling Other"
const ClassificatiomRunningSport = "Running Sport"
const ClassificatiomRunningMountain = "Running Mountain"
const ClassificatiomRunningRacing = "Running Racing"
const ClassificatiomRunningIndoor = "Running Indoor"
const ClassificatiomRunningOther = "Running Other"
const ClassificatiomWalkingSport = "Walking Sport"
const ClassificatiomWalkingMountain = "Walking Mountain"
const ClassificatiomWalkingTransport = "Walking Transport"
const ClassificatiomWalkingIndoor = "Walking Indoor"
const ClassificatiomWalkingOther = "Walking Other"
const ClassificatiomHikingSport = "Hiking Sport"
const ClassificatiomHikingMountain = "Hiking Mountain"
const ClassificatiomHikingOther = "Hiking Other"
const ClassificatiomSwimmingSport = "Swimming Sport"
const ClassificatiomSwimmingIndoor = "Swimming Indoor"
const ClassificatiomRowingSport = "Rowing Sport"
const ClassificatiomViaFerrataSport = "Via Ferrata Sport"
const ClassificatiomMotorSport = "Motor Sport"

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
