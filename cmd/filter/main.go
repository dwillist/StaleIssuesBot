package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dwillist/stale_issues/tracker"
)

func main() {
	exit(run())
}

func exit(err error) {
	if err == nil {
		os.Exit(0)
	}
	log.Printf("Error: %s\n", err)
	os.Exit(1)
}

func run() error {
	caller := tracker.Callout{}
	timer := tracker.Time{}

	trackerInstance := tracker.NewTracker(caller, timer)

	//issues, err := trackerInstance.Filter()
	//if err != nil {
	//	return err
	//}

	labels, err := trackerInstance.PostLabel()
	if err != nil {
		return err
	}

	fmt.Println(labels)

	return nil
}
