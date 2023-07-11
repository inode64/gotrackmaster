package cmd

import (
	"fmt"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var smoothGaussianDistanceCmd = &cobra.Command{
	Use:   "smoothgaussiandistance",
	Short: "Uses the Gaussian smoothing algorithm on track coordinates",
	Run: func(cmd *cobra.Command, args []string) {
		smoothGaussianDistanceExecute()
	},
}

var (
	windowSize int
	sigma      float64
)

func init() {
	rootCmd.AddCommand(smoothGaussianDistanceCmd)
	smoothGaussianDistanceCmd.Flags().IntVar(&windowSize, "windowsize", 1, "defines the window size used in the algorithm")
	smoothGaussianDistanceCmd.Flags().Float64Var(&sigma, "sigma", 1.1, "defines the sigma used in the algorithm")
}

func smoothGaussianDistanceExecute() {
	readTracks()

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		trackmaster.SmoothGaussian(g, windowSize, sigma)
		writeGPX(g, filename)
		fmt.Printf("[%v] - Smooth Gaussian distance %s\n", filename, lib.ColorRed(" (updated)"))
	}
}
