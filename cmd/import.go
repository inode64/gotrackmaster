package cmd

import (
	"fmt"
	"os"
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

func customFormat(format string, t time.Time, address geo.Address) string {
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
	return result
}

func isValidFormat(format string, t time.Time) bool {
	result := customFormat(format, t, geo.Address{Country: "Germany", CountryCode: "DE", City: "Berlin", State: "Berlin"})

	badCharMatch, _ := regexp.MatchString(`:|\\|\*|\?|"|<|>|\||\^`, format)

	return result != format && !badCharMatch
}

func isGeoAddress() bool {
	d, _ := regexp.MatchString(`\{country\}|\{countrycode\}|\{city\}|\{state\}`, directoryFormat)
	a, _ := regexp.MatchString(`\{country\}|\{countrycode\}|\{city\}|\{state\}`, archiveFormat)

	return a || d
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
		if isGeoAddress() {
			address, err = trackmaster.GetLocationStart(g)
			if err != nil {
				continue
			}
		}

		e := ImportStructure{
			source:    filename,
			directory: customFormat(directoryFormat, t, address),
			archive:   customFormat(archiveFormat, t, address),
		}
		for _, element := range importGPX {
			if element.directory == e.directory && element.archive == e.archive {
				lib.Error("Duplicate track: " + filename + " -> " + element.source)
				continue
			}
		}

		importGPX = append(importGPX, e)
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
