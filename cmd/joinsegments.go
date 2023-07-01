package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var joinSegmentsCmd = &cobra.Command{
	Use:   "joinsegments",
	Short: "Joins a segment to an adjacent segment",
	Run: func(cmd *cobra.Command, args []string) {
		joinSegmentsExecute()
	},
}

func init() {
	rootCmd.AddCommand(joinSegmentsCmd)
	joinSegmentsCmd.Flags().IntVar(&minPoints, "minpoints", 14, "Defines the minimum points of a segment to join it to the adjacent segment")
}

func joinSegmentsExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.MoveSegment(g, minPoints, true)
		writeTrack(g, filename, result)
	}
}