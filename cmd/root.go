package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

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
	lib.ReadTracks(track)
	if len(lib.Tracks) == 0 {
		lib.Error("No tracks found")
		os.Exit(1)
	}
	lib.Pass("Processing tracks...")
}

func readTrack(filename string) (gpx.GPX, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(lib.ColorYellow("Warning: GPX file could not be processed, error: ", lib.ColorRed(err)))
		return gpx.GPX{}, err
	}
	defer f.Close()

	g, err := gpx.Read(f)
	if err != nil {
		fmt.Println(lib.ColorYellow("Warning: GPX file could not be processed, error: ", lib.ColorRed(err)))
		return gpx.GPX{}, err
	}

	return *g, nil
}

func writeTrack(g gpx.GPX, filename string, result []trackmaster.GPXElementInfo) {
	if len(result) == 0 {
		fmt.Printf("[%v] - no updated need\n", filename)
	} else {
		writeGPX(g, filename)
		fmt.Printf("[%v] - Fixing %s point(s)\n", filename, lib.ColorRed(strconv.Itoa(len(result))+" (updated)"))
	}
}
