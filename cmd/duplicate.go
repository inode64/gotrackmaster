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

type duplicateStructure struct {
	startTime time.Time
	endTime   time.Time
	startLat  float64
	startLon  float64
	endLat    float64
	endLon    float64
	quality   float64
	creator   string
	filename  string
}

var (
	startDiff          int
	endDiff            int
	startDistance      int
	endDistance        int
	timeComparator     bool
	distanceComparator bool
	dup                int
	del                int
	delete             bool
)

func init() {
	rootCmd.AddCommand(duplicateCmd)
	duplicateCmd.Flags().IntVar(&startDiff, "startdiff", 0, "Time in seconds from the beginning of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDiff, "enddiff", 0, "Time in seconds from the end of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&startDistance, "startDistance", 0, "Distance in meters from the beginning position of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().IntVar(&endDistance, "endDistance", 0, "Distance in meters from the end position of the track to determine if they are duplicates (set 0 to not use this rule)")
	duplicateCmd.Flags().BoolVar(&timeComparator, "timeComparator", false, "Takes time start and end to determine if they are duplicates")
	duplicateCmd.Flags().BoolVar(&distanceComparator, "distanceComparator", false, "Requires distance start and end to determine if they are duplicates")
	duplicateCmd.Flags().BoolVar(&delete, "delete", false, "Delete duplicate only when equal creator and quality of track")
}

func checkTime(t, d time.Time, sec int) bool {
	return t.After(d.Add(time.Duration(-sec)*time.Second)) && t.Before(d.Add(time.Duration(sec)*time.Second))
}

func checkPosition(lat1, lon1, lat2, lon2 float64, distance int) bool {
	return trackmaster.HaversineDistance(lat1, lon1, lat2, lon2) < float64(distance)
}

func showDuplicate(d, new duplicateStructure, status string) bool {
	lib.Error(fmt.Sprintf("Duplicate found: %v [%v]", showNameTrack(d.filename, d.creator, d.quality), status))
	dup++
	if delete && d.creator == new.creator && d.quality == new.quality {
		del++
		lib.Info(fmt.Sprintf("Deleting %v", d.filename))
		if !dryRun {
			os.Remove(d.filename)
			return true
		}
	}
	return false
}

func showNameTrack(filename, creator string, quality float64) string {
	return fmt.Sprintf("%v (%s/%0.0f)", filename, creator, quality)
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
	if startDiff == 0 && endDiff == 0 && startDistance == 0 && endDistance == 0 {
		lib.Error("You must specify at least one rule")
		os.Exit(1)
	}

	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		lib.Error(err.Error())
		os.Exit(1)
	}

	readTracks()

	var duplicateGPX []duplicateStructure

	for _, filename := range lib.Tracks {
		g, err := readTrack(filename)
		if err != nil {
			continue
		}

		quality := trackmaster.QualityTrack(g)
		creator := trackmaster.GetCreator(g)

		fmt.Printf("Getting info from: %v\n", showNameTrack(filename, creator, quality))

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

		new := duplicateStructure{
			startTime: ts,
			endTime:   te,
			startLat:  ps.Lat,
			startLon:  ps.Lon,
			endLat:    pe.Lat,
			endLon:    pe.Lon,
			quality:   quality,
			creator:   creator,
			filename:  filename,
		}

		// check if start time is same other track with margin of startDiff
		if startDiff != 0 {
			for _, d := range duplicateGPX {
				if checkTime(ts, d.startTime, startDiff) {
					if timeComparator && endDiff != 0 && checkTime(te, d.endTime, endDiff) {
						if showDuplicate(d, new, "start and end time") {
							goto next
						}
					} else {
						if showDuplicate(d, new, "start time") {
							goto next
						}
					}
				}
			}
		} else if endDiff != 0 {
			for _, d := range duplicateGPX {
				if checkTime(te, d.endTime, endDiff) {
					if showDuplicate(d, new, "end time") {
						goto next
					}
				}
			}
		}

		if startDistance != 0 {
			for _, d := range duplicateGPX {
				if checkPosition(ps.Lat, ps.Lon, d.startLat, d.startLon, startDistance) {
					if distanceComparator && endDistance != 0 && checkPosition(pe.Lat, pe.Lon, d.endLat, d.endLon, endDistance) {
						if showDuplicate(d, new, "start and end position") {
							goto next
						}
					} else {
						if showDuplicate(d, new, "start position") {
							goto next
						}
					}
				}
			}
		} else if endDistance != 0 {
			for _, d := range duplicateGPX {
				if checkPosition(pe.Lat, pe.Lon, d.endLat, d.endLon, endDistance) {
					if showDuplicate(d, new, "end position") {
						goto next
					}
				}
			}
		}

		duplicateGPX = append(duplicateGPX, new)
	next:
	}
	lib.Pass(fmt.Sprintf("Found %d duplicate tracks", dup))
	lib.Pass(fmt.Sprintf("Deleted %d duplicate tracks", del))
}
