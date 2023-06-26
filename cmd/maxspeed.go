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

var maxSpeedCmd = &cobra.Command{
	Use:   "maxspeed",
	Short: "Remove points using max speed",
	Run: func(cmd *cobra.Command, args []string) {
		elevationExecute()
	},
}
var maxSpeed float64

func init() {
	rootCmd.AddCommand(maxSpeedCmd)
	rootCmd.PersistentFlags().Float64Var(&maxSpeed, "maxspeed", 200.0, "set the maximum speed to remove from track")
}

func maxSpeedExecute() {
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

		result := trackmaster.MaxSpeed(*g, maxSpeed, true)
		if len(result) > 0 {
			writeGPX(*g, filename)
			fmt.Printf("[%v] - Fixing %s point(s)\n", filename, lib.ColorRed(strconv.Itoa(len(result))+" (updated)"))
		}
	}
}
