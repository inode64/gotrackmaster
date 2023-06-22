package cmd

import (
	"fmt"
	"os"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-gpx"
)

var elevationCmd = &cobra.Command{
	Use:   "elevation",
	Short: "Update elevation using SRTM data",
	Run: func(cmd *cobra.Command, args []string) {
		elevationExecute()
	},
}
var accuracy int16

func init() {
	rootCmd.AddCommand(elevationCmd)
	rootCmd.PersistentFlags().Int16Var(&accuracy, "accuracy", 60, "set the minimum accuracy to update the elevation")
}

func elevationExecute() {
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

		num, err := trackmaster.ElevationSRTMAccuracy(*g)
		if err != nil {
			fmt.Println(lib.ColorYellow("Warning: Elevation SRTM could not be processed, error: ", lib.ColorRed(err)))
			continue
		}
		if int16(num) <= accuracy {
			fmt.Printf("[%v] - %s\n", filename, lib.ColorGreen("updated"))
		}

		fmt.Printf("[%v] - Accuracy %s\n", filename, lib.ColorGreen(num))
	}
}
