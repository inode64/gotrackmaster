package cmd

import (
	"fmt"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Show the quality of track",
	Run: func(cmd *cobra.Command, args []string) {
		qualityExecute()
	},
}

func init() {
	rootCmd.AddCommand(qualityCmd)
}

func qualityExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}
		quality := trackmaster.QualityTrack(g)
		fmt.Printf("[%v] - %s\n", filename, lib.ColorGreen(quality))
	}
}
