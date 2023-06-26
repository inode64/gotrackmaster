package cmd

import (
	"fmt"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var classificationCmd = &cobra.Command{
	Use:   "classification",
	Short: "Classify a track according to the type of activity",
	Run: func(cmd *cobra.Command, args []string) {
		classificationExecute()
	},
}

func init() {
	rootCmd.AddCommand(classificationCmd)
}

func classificationExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		kind := trackmaster.ClassificationTrack(filename)
		fmt.Printf("[%v] - %s\n", filename, lib.ColorGreen(kind))
	}
}
