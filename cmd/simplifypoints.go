package cmd

import (
	"math"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var simplifyPointsCmd = &cobra.Command{
	Use:   "simplifypoints",
	Short: "Simplify the track by removing very close points",
	Run: func(cmd *cobra.Command, args []string) {
		simplifyPointsExecute()
	},
}
var distance float64

func init() {
	rootCmd.AddCommand(simplifyPointsCmd)
	simplifyPointsCmd.Flags().Float64Var(&distance, "distance", 0.5, "set minimum distance of the points to join them")
}

func simplifyPointsExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveStops(g, 0.0, distance, math.MaxFloat64, 0, true)

		writeTrack(g, filename, result)
	}
}
