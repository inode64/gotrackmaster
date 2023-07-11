package lib

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	ColorBlue   = color.New(color.FgBlue).SprintFunc()
	ColorGreen  = color.New(color.FgGreen).SprintFunc()
	ColorRed    = color.New(color.FgRed).SprintFunc()
	ColorYellow = color.New(color.FgYellow).SprintFunc()
)

func Error(s string) {
	logrus.Error(Colorize(s, color.FgRed))
}

func Warning(s string) {
	logrus.Warn(Colorize(s, color.FgRed))
}

func Notice(s string) {
	logrus.Info(Colorize(s, color.FgYellow))
}

func Info(s string) {
	logrus.Info(Colorize(s, color.FgGreen))
}

func Pass(s string) {
	logrus.Info(Colorize(s, color.FgBlue))
}

func Colorize(s string, c color.Attribute) string {
	return color.New(c).SprintFunc()(s)
}
