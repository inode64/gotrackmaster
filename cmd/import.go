package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/codingsince1985/geo-golang"
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/ringsaturn/tzf"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports tracks and sorts them into a new directory structure",
	Run: func(cmd *cobra.Command, args []string) {
		importExecute()
	},
}

type ImportStructure struct {
	source    string
	directory string
	archive   string
}

var (
	destination     string
	directoryFormat string
	archiveFormat   string
	finder          tzf.F
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVar(&destination, "destination", "", "destination directory to classify the tracks")
	importCmd.Flags().StringVar(&directoryFormat, "directoryformat", "", "directory format for the tracks")
	importCmd.Flags().StringVar(&archiveFormat, "archiveformat", "", "archive format for the tracks")
}

func customFormat(format string, t time.Time, address geo.Address, degree1 string, degree5 string, original string, kind string, creator string) string {
	result := format
	result = strings.ReplaceAll(result, "{year}", fmt.Sprintf("%d", t.Year()))
	result = strings.ReplaceAll(result, "{month}", fmt.Sprintf("%02d", t.Month()))
	result = strings.ReplaceAll(result, "{day}", fmt.Sprintf("%02d", t.Day()))
	result = strings.ReplaceAll(result, "{hour}", fmt.Sprintf("%02d", t.Hour()))
	result = strings.ReplaceAll(result, "{minute}", fmt.Sprintf("%02d", t.Minute()))
	result = strings.ReplaceAll(result, "{country}", address.Country)
	result = strings.ReplaceAll(result, "{countrycode}", address.CountryCode)
	result = strings.ReplaceAll(result, "{city}", address.City)
	result = strings.ReplaceAll(result, "{state}", address.State)
	result = strings.ReplaceAll(result, "{degree1}", degree1)
	result = strings.ReplaceAll(result, "{degree0.5}", degree5)
	result = strings.ReplaceAll(result, "{original}", original)
	result = strings.ReplaceAll(result, "{kind}", kind)
	result = strings.ReplaceAll(result, "{creator}", creator)
	return result
}

func isValidFormat(format string, t time.Time) bool {
	result := customFormat(format, t, geo.Address{Country: "Germany", CountryCode: "DE", City: "Berlin", State: "Berlin"}, "0", "0", "original", "trail running", "Strava")

	badCharMatch, _ := regexp.MatchString(`:|\\|\*|\?|"|<|>|\||\^`, format)

	return result != format && !badCharMatch
}

func isGeoAddress() bool {
	d, _ := regexp.MatchString(`\{country\}|\{countrycode\}|\{city\}|\{state\}`, directoryFormat)
	a, _ := regexp.MatchString(`\{country\}|\{countrycode\}|\{city\}|\{state\}`, archiveFormat)

	return a || d
}

func isDegree1() bool {
	d, _ := regexp.MatchString(`\{degree1\}`, directoryFormat)
	a, _ := regexp.MatchString(`\{degree1\}`, archiveFormat)

	return a || d
}

func isDegree5() bool {
	d, _ := regexp.MatchString(`\{degree0.5\}`, directoryFormat)
	a, _ := regexp.MatchString(`\{degree0.5\}`, archiveFormat)

	return a || d
}

func appendTrack(filename string, t time.Time, address geo.Address, gpx []ImportStructure, degree1 string, degree5 string, creator string) []ImportStructure {
	file := filepath.Base(filename)
	extension := filepath.Ext(file)
	name := file[:len(file)-len(extension)]
	kind := trackmaster.ClassificationTrack(filename)

	e := ImportStructure{
		source:    filename,
		directory: customFormat(directoryFormat, t, address, degree1, degree5, name, kind, creator),
		archive:   customFormat(archiveFormat, t, address, degree1, degree5, name, kind, creator),
	}
	for _, element := range gpx {
		if element.directory == e.directory && element.archive == e.archive {
			lib.Error("Duplicate track: " + filename + " == " + element.source)
			continue
		}
	}

	return append(gpx, e)
}

func importExecute() {
	if destination == "" {
		lib.Error("Destination directory is missing")
		os.Exit(1)
	}
	if directoryFormat != "" && !isValidFormat(directoryFormat, time.Now()) {
		lib.Error("Directory format is wrong")
		os.Exit(1)
	}
	if !isValidFormat(archiveFormat, time.Now()) {
		lib.Error("Archive format is wrong")
		os.Exit(1)
	}

	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		lib.Error(err.Error())
		os.Exit(1)
	}

	readTracks()

	var importGPX []ImportStructure
	var address geo.Address

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		fmt.Printf("Getting info from: %v\n", filename)

		t := trackmaster.GetTimeStart(g, finder)
		if t.IsZero() {
			continue
		}
		creator := trackmaster.GetCreator(g)

		if isGeoAddress() {
			address, err = trackmaster.GetLocationStart(g)
		}
		if isDegree1() || isDegree5() {
			bounds := trackmaster.GetBounds(g)
			if trackmaster.IsBoundsValid(bounds) {
				degree1 := trackmaster.CalculateTiles(bounds, 1)
				degree5 := trackmaster.CalculateTiles(bounds, 0.5)
				if isDegree1() {
					for _, element1 := range degree1 {
						if isDegree5() {
							for _, element5 := range degree5 {
								importGPX = appendTrack(filename, t, address, importGPX, element1, element5, creator)
							}
						} else {
							importGPX = appendTrack(filename, t, address, importGPX, element1, "", creator)
						}
					}
				} else {
					for _, element5 := range degree5 {
						importGPX = appendTrack(filename, t, address, importGPX, "", element5, creator)
					}
				}
			} else {
				importGPX = appendTrack(filename, t, address, importGPX, "", "", creator)
			}
		} else {
			importGPX = appendTrack(filename, t, address, importGPX, "", "", creator)
		}
	}

	lib.Pass("Moving tracks...")

	for _, element := range importGPX {
		g, err := readTrack(element.source)
		if err != nil {
			continue
		}

		target := destination + "/" + element.directory + "/" + element.archive + ".gpx"
		fmt.Printf("[%v] -> %v\n", element.source, target)

		if !dryRun {
			err = os.MkdirAll(destination+"/"+element.directory, os.ModePerm)
			if err != nil {
				lib.Error(err.Error())
				os.Exit(1)
			}

			writeGPX(g, target)
		}
	}
}
