package cmd

import (
	"fmt"
	"os"

	"github.com/inode64/gotrackmaster/lib"
	trackmaster "github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var classificationCmd = &cobra.Command{
	Use:   "classification",
	Short: "Synchronize Media Data from track GPX",
	Long:  `Using a gpx track, analyze a directory with images or movies and add the GPS positions`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		MExecute()
	},
}

func init() {
	rootCmd.AddCommand(classificationCmd)
}

func MExecute() {
	lib.ReadTracks(track, true)
	lib.Pass("Processing tracks...")

	if len(lib.Tracks) == 0 {
		os.Exit(1)
	}

	for _, filename := range lib.Tracks {
		kind := trackmaster.ClassificationTrack(filename)
		fmt.Printf("[%v] - %s\n", filename, lib.ColorGreen(kind))
	}
}
