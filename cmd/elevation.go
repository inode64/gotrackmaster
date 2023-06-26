package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
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
	readTracks()

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
		if int16(num) > accuracy {
			fmt.Printf("[%v] - Accuracy %s\n", filename, lib.ColorGreen(num))
		} else {
			if !dryRun {
				trackmaster.ElevationSRTM(*g)

				writeGPX(*g, filename)
			}
			fmt.Printf("[%v] - Accuracy %s\n", filename, lib.ColorRed(strconv.Itoa(num)+" (updated)"))
		}
	}
}
