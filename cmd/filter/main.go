package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dwillist/stale_issues/tracker"
	"github.com/pkg/errors"

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

	issues, err := trackerInstance.Filter()
	if err != nil {
		return err
	}

	_, success, err := trackerInstance.PostLabel()
	if err != nil {
		return err
	} else if success {
		fmt.Println("Created new label for stale issues")
	} else {
		fmt.Println("Stale issues label already exists")
	}

	for _, story  := range issues {
		updatedStory, success, err := trackerInstance.UpdateStory(story)
		if err != nil {
			return errors.Wrap(err, "failed to apply Stale Issue")
		} else if success {
			fmt.Printf("Updated story name: %s  id:%s", updatedStory.Name, updatedStory.ID)
		} else {
			fmt.Printf("Failed to update name: %s id:%s", updatedStory.Name, updatedStory.ID)
		}
	}

	return nil
}
