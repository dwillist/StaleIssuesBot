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

	trackerInstance := tracker.NewTracker(caller)

	staleIssues, err := trackerInstance.Search()
	if err != nil {
		return err
	}

	fmt.Printf("Stale Issue: %v", staleIssues)

	return nil
}
