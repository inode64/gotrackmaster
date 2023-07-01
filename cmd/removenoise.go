package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var removeNoiseCmd = &cobra.Command{
	Use:   "removenoise",
	Short: "Remove intense noise on the track",
	Run: func(cmd *cobra.Command, args []string) {
		removeNoiseExecute()
	},
}

var maxPoints int

func init() {
	rootCmd.AddCommand(removeNoiseCmd)
	removeNoiseCmd.Flags().Float64Var(&maxDistance, "maxdistance", 6.0, "set the maximum distance allowed within a stop")
	removeNoiseCmd.Flags().Float64Var(&maxElevation, "maxelevation", 1.1, "set the maximum lift allowed within a stop")
	removeNoiseCmd.Flags().IntVar(&maxPoints, "maxpoints", 4, "set the maximum amount of points")
}

func removeNoiseExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveNoise(g, maxDistance, maxElevation, maxPoints, true)
		writeTrack(g, filename, result)
	}
}
