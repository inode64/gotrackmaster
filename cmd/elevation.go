package cmd

import (
	"fmt"
	"strconv"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
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
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		num, err := trackmaster.ElevationSRTMAccuracy(g)
		if err != nil {
			fmt.Println(lib.ColorYellow("Warning: Elevation SRTM could not be processed, error: ", lib.ColorRed(err)))
			continue
		}
		if int16(num) > accuracy {
			fmt.Printf("[%v] - Accuracy %s\n", filename, lib.ColorGreen(num))
		} else {
			if !dryRun {
				trackmaster.ElevationSRTM(g)

				writeGPX(g, filename)
			}
			fmt.Printf("[%v] - Accuracy %s\n", filename, lib.ColorRed(strconv.Itoa(num)+" (updated)"))
		}
	}
}
