package lib

import (
	"fmt"
	"log"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/karrick/godirwalk"
	"github.com/twpayne/go-gpx"
)

var (
	Tracks     []string
	TrackValid int
	TrackError int
)

func ReadGPX(filename string, valid bool) {
	mtype, err := mimetype.DetectFile(filename)
	if err != nil {
		log.Fatal(ColorRed(err))
	}

	if !mtype.Is("application/gpx+xml") && !mtype.Is("text/xml") {
		return
	}

	fmt.Printf("Reading: %v \n", filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(ColorYellow("Warning: GPX file could not be processed, error: ", ColorRed(err)))
		TrackError++
		return
	}
	defer file.Close()

	_, err = gpx.Read(file)
	if err != nil {
		fmt.Println(ColorYellow("Warning: GPX file could not be processed, error: ", ColorRed(err)))
		TrackError++
		return
	}

	Tracks = append(Tracks, filename)
	TrackValid++
}

func ReadGPXDir(trackDir string, valid bool) {
	err := godirwalk.Walk(trackDir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				return nil // do not remove directory that was provided top-level directory
			}

			ReadGPX(path, valid)

			return nil
		},
		Unsorted: false,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ReadTracks(track string, valid bool) {
	fileInfo, err := os.Stat(track)
	if err != nil {
		log.Fatal(ColorRed("No open GPX path"))
	}

	Pass("Reading tracks...")

	if fileInfo.IsDir() {
		ReadGPXDir(track, valid)
	} else {
		ReadGPX(track, valid)
	}

	if len(Tracks) == 0 {
		Warning("There is no track processed")
	}

	if TrackError == 0 {
		fmt.Printf(ColorGreen("Located %d track(s)"), TrackValid)
	} else {
		fmt.Printf(ColorYellow("Located %d track(s), %d with error(s)"), TrackValid, TrackError)
	}
}
