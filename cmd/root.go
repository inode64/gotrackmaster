package cmd

import (
	"github.com/spf13/cobra"
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
	Args:    cobra.MinimumNArgs(0),
	Version: "1.0.0",
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
