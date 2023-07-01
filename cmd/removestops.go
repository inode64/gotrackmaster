package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var removeStopsCmd = &cobra.Command{
	Use:   "removestops",
	Short: "Remove stops on a track",
	Run: func(cmd *cobra.Command, args []string) {
		removeStopsExecute()
	},
}

var minSeconds float64
var minPoints int
var maxElevation float64
var maxDistance float64

func init() {
	rootCmd.AddCommand(removeStopsCmd)
	removeStopsCmd.Flags().Float64Var(&maxDistance, "maxdistance", 5.0, "set the maximum distance allowed within a stop")
	removeStopsCmd.Flags().Float64Var(&minSeconds, "minseconds", 90.0, "set the minimum time that is considered a stop")
	removeStopsCmd.Flags().Float64Var(&maxElevation, "maxelevation", 0.5, "set the maximum lift allowed within a stop")
	removeStopsCmd.Flags().IntVar(&minPoints, "minpoints", 3.0, "set the minimum amount of points")
}

func removeStopsExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveStops(g, minSeconds, maxDistance, maxElevation, minPoints, true)
		writeTrack(g, filename, result)
	}
}
