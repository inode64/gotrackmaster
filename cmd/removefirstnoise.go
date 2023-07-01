package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var removeFirstNoiseCmd = &cobra.Command{
	Use:   "removefirstnoise",
	Short: "Removes noise from the first generated points of the track due to poor signal quality",
	Run: func(cmd *cobra.Command, args []string) {
		removeFirstNoiseExecute()
	},
}

func init() {
	rootCmd.AddCommand(removeFirstNoiseCmd)
}

func removeFirstNoiseExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveFirstNoise(g, true)

		writeTrack(g, filename, result)
	}
}
