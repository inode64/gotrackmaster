package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/inode64/gotrackmaster/lib"
	"github.com/inode64/gotrackmaster/trackmaster"
	"github.com/ringsaturn/tzf"
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
	timeComparator     bool
	distanceComparator bool
)

func init() {
	rootCmd.AddCommand(duplicateCmd)
	duplicateCmd.Flags().IntVar(&startDiff, "startdiff", 0, "Time in seconds from the beginning of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDiff, "enddiff", 0, "Time in seconds from the end of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&startDistance, "startDistance", 0, "Distance in meters from the beginning position of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDistance, "endDistance", 0, "Distance in meters from the end position of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().BoolVar(&timeComparator, "timeComparator", false, "Takes time start and end to determine if they are duplicates")
	duplicateCmd.Flags().BoolVar(&distanceComparator, "distanceComparator", false, "Requires distance start and end to determine if they are duplicates")
}

func checkTime(t, d time.Time, sec int) bool {
	return t.After(d.Add(time.Duration(-sec)*time.Second)) && t.Before(d.Add(time.Duration(sec)*time.Second))
}

func CheckPosition(lat1, lon1, lat2, lon2 float64, distance int) bool {
	return trackmaster.HaversineDistance(lat1, lon1, lat2, lon2) < float64(distance)
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

	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		lib.Error(err.Error())
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

		ts := trackmaster.GetTimeStart(g, finder)
		te := trackmaster.GetTimeEnd(g, finder)

		// only add if start and end time are valid
		if (ts.IsZero() || te.IsZero()) && startDiff != 0 && endDiff != 0 && startDistance == 0 && endDistance == 0 {
			continue
		}

		ps := trackmaster.GetPositionStart(g)
		pe := trackmaster.GetPositionEnd(g)
		// without position start and end we can't calculate distance
		if (ps.Lat == 0 && ps.Lon == 0) || (pe.Lat == 0 && pe.Lon == 0) {
			continue
		}

		// check if start time is same other track with margin of startDiff
		if startDiff != 0 {
			for _, d := range duplicateGPX {
				if checkTime(ts, d.startTime, startDiff) {
					if timeComparator && endDiff != 0 && checkTime(te, d.endTime, endDiff) {
						fmt.Printf("Duplicate found: %v [start and end time]\n", filename)
						continue
					} else {
						fmt.Printf("Duplicate found: %v [start time]\n", filename)
						continue
					}
				}
			}
		} else if endDiff != 0 {
			for _, d := range duplicateGPX {
				if checkTime(te, d.endTime, endDiff) {
					fmt.Printf("Duplicate found: %v [end time]\\n", filename)
					continue
				}
			}
		}

		if startDistance != 0 {
			for _, d := range duplicateGPX {
				if CheckPosition(ps.Lat, ps.Lon, d.startLat, d.startLon, startDistance) {
					if distanceComparator && endDistance != 0 && CheckPosition(pe.Lat, pe.Lon, d.endLat, d.endLon, endDistance) {
						fmt.Printf("Duplicate found: %v [start and end position]\n", filename)
						continue
					} else {
						fmt.Printf("Duplicate found: %v [start position]\n", filename)
						continue
					}
				}
			}
		} else if endDistance != 0 {
			for _, d := range duplicateGPX {
				if CheckPosition(pe.Lat, pe.Lon, d.endLat, d.endLon, endDistance) {
					fmt.Printf("Duplicate found: %v [end position]\n", filename)
					continue
				}
			}
		}

		duplicateGPX = append(duplicateGPX, DuplicateStructure{
			startTime: ts,
			endTime:   te,
			startLat:  ps.Lat,
			startLon:  ps.Lon,
			endLat:    pe.Lat,
			endLon:    pe.Lon,
		})
	}
}
