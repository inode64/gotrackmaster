package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var removeLastMaxSpeedCmd = &cobra.Command{
	Use:   "removelastmaxspeed",
	Short: "Removes from the end of the track when you do not stop recording and get into a vehicle",
	Run: func(cmd *cobra.Command, args []string) {
		removeLastMaxSpeedExecute()
	},
}

func init() {
	rootCmd.AddCommand(removeLastMaxSpeedCmd)
	removeLastMaxSpeedCmd.Flags().Float64Var(&maxSpeed, "maxspeed", 14.0, "set the maximum speed to remove from the end of the track")
}

func removeLastMaxSpeedExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveLastMaxSpeed(g, maxSpeed, true)

		writeTrack(g, filename, result)
	}
}
