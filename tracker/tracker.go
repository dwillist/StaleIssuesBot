package tracker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwillist/stale_issues/resources"
)

const (
	Endpoint         = "https://www.pivotaltracker.com/services/v5/projects/1042066/search?query=label%3Agithub-issue%20AND%20-state%3Aaccepted%20-state%3Afinished%20-state%3Adelivered"
	StaleAfterMonths = 1
)

type Tracker struct {
	Caller Caller
	Timer  Timer
}

type Caller interface {
	Call(endpoint string) (string, error)
}

type Timer interface {
	Time() time.Time
}

func NewTracker(caller Caller, timer Timer) Tracker {
	return Tracker{
		Caller: caller,
		Timer:  timer,
	}
}

func (t Tracker) Search() ([]resources.Story, error) {
	var responseStruct resources.Search
	var result []resources.Story

	response, err := t.Caller.Call(Endpoint)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(response), &responseStruct); err != nil {
		return result, err
	}

	result = responseStruct.Stories.Stories

	return result, nil
}

func (t Tracker) Filter() ([]resources.Story, error) {
	var result []resources.Story

	issues, err := t.Search()
	if err != nil {
		return result, err
	}
	result = t.filterIssues(issues)

	return result, nil
}

func (t Tracker) filterIssues(stories []resources.Story) []resources.Story {
	var result []resources.Story

	fmt.Println("total issues:", len(stories))

	for _, story := range stories {
		if t.isStale(story) {
			result = append(result, story)
		}
	}

	fmt.Println("Stale Issues:", len(result))

	return result
}

func (t Tracker) isStale(story resources.Story) bool {
	return !story.UpdatedAt.AddDate(0, StaleAfterMonths, 0).After(t.Timer.Time())
}
