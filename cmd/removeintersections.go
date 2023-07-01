package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var removeIntersectionsCmd = &cobra.Command{
	Use:   "removeintersections",
	Short: "Remove intersections on the track",
	Run: func(cmd *cobra.Command, args []string) {
		removeIntersectionsExecute()
	},
}

func init() {
	rootCmd.AddCommand(removeIntersectionsCmd)
	removeIntersectionsCmd.Flags().IntVar(&maxPoints, "maxpoints", 6, "set the maximum amount of points")
}

func removeIntersectionsExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.RemoveIntersections(g, maxPoints, true)
		writeTrack(g, filename, result)
	}
}
