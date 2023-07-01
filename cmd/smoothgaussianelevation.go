package cmd

import (
	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var smoothGaussianElevationCmd = &cobra.Command{
	Use:   "smoothgaussianelevation",
	Short: "Uses the Gaussian smoothing algorithm to adjust elevation",
	Run: func(cmd *cobra.Command, args []string) {
		smoothGaussianElevationExecute()
	},
}

func init() {
	rootCmd.AddCommand(smoothGaussianElevationCmd)
	smoothGaussianElevationCmd.Flags().Float64Var(&maxElevation, "maxelevation", 1.5, "defines the maximum vertical speed to perform a smoothing")
}

func smoothGaussianElevationExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		result := trackmaster.MaxSpeedVertical(g, maxElevation, true)
		writeTrack(g, filename, result)
	}
}
