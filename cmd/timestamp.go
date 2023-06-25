package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-gpx"
)

var timeCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Update timestamp in all GPX file",
	Long:  "Corrects all the timestamps that are missing or those that are outside the timeline in the track",
	Run: func(cmd *cobra.Command, args []string) {
		timExecute()
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}

func timExecute() {
	if verbose {
		trackmaster.Log.SetLevel(logrus.DebugLevel)
	}

	lib.ReadTracks(track, true)
	lib.Pass("Processing tracks...")

	if len(lib.Tracks) == 0 {
		os.Exit(1)
	}

	for _, filename := range lib.Tracks {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(lib.ColorYellow("Warning: GPX file could not be processed, error: ", lib.ColorRed(err)))
			continue
		}
		defer f.Close()

		g, err := gpx.Read(f)
		if err != nil {
			fmt.Println(lib.ColorYellow("Warning: GPX file could not be processed, error: ", lib.ColorRed(err)))
			continue
		}

		if trackmaster.TimeEmpty(*g) {
			fmt.Println(lib.ColorRed("Error: GPX file hasn't any time"))
			continue
		}

		quality := trackmaster.TimeQuality(*g)
		if quality == 100 {
			fmt.Printf("[%v] - Tack with all correct timestamp \n", filename)
			continue
		}
		if quality == -1 {
			fmt.Println(lib.ColorRed("Error: GPX file empty"))
			continue
		}

		num := trackmaster.FixTimesTrack(*g, true)
		quality = trackmaster.TimeQuality(*g)
		if quality != 100 {
			fmt.Printf("[%v] - Timestamp that could not be corrected\n", filename)
		} else {
			fmt.Printf("[%v] - Timestamp that have been fixed %s\n", filename, lib.ColorRed(strconv.Itoa(num)+" (updated)"))
			if !dryRun {
				writeGPX(*g, filename)
			}
		}
	}
}
