package cmd

import (
	"encoding/xml"
	"os"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-gpx"
)

var (
	dryRun  bool
	force   bool
	verbose bool
	track   string
)

var rootCmd = &cobra.Command{
	Use:   "gotrackmaster",
	Short: "Manage GPX tracks",
	Long: `A versatile Go-based toolkit for comprehensive GPX track analysis and optimization.
Features include maximum speed calculations, slope computations, removal of erratic points,
track simplification, and more. Ideal for outdoor enthusiasts, athletes,
and GIS professionals seeking insights from their GPX data.`,
	Version: "1.0.0",
	Args:    cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Performs the actions without writing to the files")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "Force update even overwriting previous GPS data")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show more information")
	rootCmd.PersistentFlags().StringVar(&track, "track", "", "GPX track or a directory of GPX tracks")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func writeGPX(g gpx.GPX, filename string) {
	if dryRun {
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		lib.Error(err.Error())
		return
	}
	defer f.Close()

	// write xml header
	_, err = f.WriteString(xml.Header)
	if err != nil {
		lib.Error(err.Error())
		return
	}

	if err := g.WriteIndent(f, "", "  "); err != nil {
		lib.Error(err.Error())
	}
}

func readTracks() {
	if verbose {
		trackmaster.Log.SetLevel(logrus.DebugLevel)
	}

	lib.ReadTracks(track, true)
	lib.Pass("Processing tracks...")

	if len(lib.Tracks) == 0 {
		os.Exit(1)
	}
}
