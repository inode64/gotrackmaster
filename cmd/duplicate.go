package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/spf13/cobra"
)

var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "Search for duplicate tracks by time (start/end) and/or position (start/end)",
	Run: func(cmd *cobra.Command, args []string) {
		duplicateExecute()
	},
}

type DuplicateStructure struct {
	startTime time.Time
	endTime   time.Time
	startLat  float64
	startLon  float64
	endLat    float64
	endLon    float64
}

var (
	startDiff          int
	endDiff            int
	startDistance      int
	endDistance        int
	comparator         bool
	timeComparator     bool
	distanceComparator bool
)

func init() {
	rootCmd.AddCommand(duplicateCmd)
	duplicateCmd.Flags().IntVar(&startDiff, "startdiff", 0, "Time in seconds from the beginning of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDiff, "enddiff", 0, "Time in seconds from the end of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDiff, "startDistance", 0, "Distance in meters from the beginning of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDiff, "endDistance", 0, "Distance in meters from the end of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().BoolVar(&comparator, "comparator", false, "Requires time and distance to determine if they are duplicates")
	duplicateCmd.Flags().BoolVar(&timeComparator, "timeComparator", false, "Takes time start and end to determine if they are duplicates")
	duplicateCmd.Flags().BoolVar(&distanceComparator, "distanceComparator", false, "Requires distance start and end to determine if they are duplicates")
}

func duplicateExecute() {
	if startDiff < 0 {
		lib.Error("Start diff must be positive")
		os.Exit(1)
	}
	if endDiff < 0 {
		lib.Error("End diff must be positive")
		os.Exit(1)
	}
	if startDistance < 0 {
		lib.Error("Start distance must be positive")
		os.Exit(1)
	}
	if endDistance < 0 {
		lib.Error("End distance must be positive")
		os.Exit(1)
	}

	readTracks()

	var duplicateGPX []DuplicateStructure

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		fmt.Printf("Getting info from: %v\n", filename)

		t := trackmaster.GetTimeStart(g, finder)
		if t.IsZero() {
			continue
		}
		creator := trackmaster.GetCreator(g)

		if isGeoAddress() {
			address, err = trackmaster.GetLocationStart(g)
		}
		if isDegree1() || isDegree5() {
			bounds := trackmaster.GetBounds(g)
			if trackmaster.IsBoundsValid(bounds) {
				degree1 := trackmaster.CalculateTiles(bounds, 1)
				degree5 := trackmaster.CalculateTiles(bounds, 0.5)
				if isDegree1() {
					for _, element1 := range degree1 {
						if isDegree5() {
							for _, element5 := range degree5 {
								importGPX = appendTrack(filename, t, address, importGPX, element1, element5, creator)
							}
						} else {
							importGPX = appendTrack(filename, t, address, importGPX, element1, "", creator)
						}
					}
				} else {
					for _, element5 := range degree5 {
						importGPX = appendTrack(filename, t, address, importGPX, "", element5, creator)
					}
				}
			} else {
				importGPX = appendTrack(filename, t, address, importGPX, "", "", creator)
			}
		} else {
			importGPX = appendTrack(filename, t, address, importGPX, "", "", creator)
		}
	}
}
