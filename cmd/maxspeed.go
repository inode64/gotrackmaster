package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var maxSpeedCmd = &cobra.Command{
	Use:   "maxspeed",
	Short: "Remove points using max speed",
	Run: func(cmd *cobra.Command, args []string) {
		maxSpeedExecute()
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
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.MaxSpeed(g, maxSpeed, true)
		writeTrack(g, filename, result)
	}
}
