package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var lostElevationCmd = &cobra.Command{
	Use:   "lostelevation",
	Short: "Fixes elevation when elevation changes abruptly",
	Run: func(cmd *cobra.Command, args []string) {
		lostElevationExecute()
	},
}

func init() {
	rootCmd.AddCommand(lostElevationCmd)
}

func lostElevationExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.LostElevation(g, true)
		writeTrack(g, filename, result)
	}
}
