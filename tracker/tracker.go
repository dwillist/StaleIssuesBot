package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dwillist/stale_issues/resources"
)

const (
	ProjectID        = "1042066"
	SearchEndpoint   = "https://www.pivotaltracker.com/services/v5/projects/" + ProjectID + "/search?query=label%3Agithub-issue%20AND%20-state%3Aaccepted%20-state%3Afinished%20-state%3Adelivered"
	LabelsEndpoint   = "https://www.pivotaltracker.com/services/v5/projects/" + ProjectID + "/labels"
	StaleAfterMonths = 1
	StaleLabel       = "stale-issue"
)

type Tracker struct {
	Caller Caller
	Timer  Timer
}

type Caller interface {
	Get(endpoint string) ([]byte, error)
	Post(endpoint string, data []byte) ([]byte, error)
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

	response, err := t.Caller.Get(SearchEndpoint)
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

// What do we want this to do ????
// Post a label successfully or

// should probably return the label object at some point here
// both cases should give us such
func (t Tracker) PostLabel() (resources.Label, bool, error) {
	newLabel := resources.Label{Name: StaleLabel}
	labelBytes, err := json.Marshal(newLabel)
	if err != nil {
		return resources.Label{}, false, err
	}
	postResponse, err := t.Caller.Post(LabelsEndpoint, labelBytes)
	if err != nil {
		return resources.Label{}, false, err
	}

	var errorResponse resources.TrackerError
	if err := json.Unmarshal(postResponse, &errorResponse); err == nil {
		// need to do some work in here...
		label, err := t.getLabelFromName(StaleLabel)
		if err != nil {
			return resources.Label{}, false, err
		}
		return label, false, nil

	}

	var successResponse resources.Label
	if err := json.Unmarshal(postResponse, &successResponse); err == nil {
		return successResponse, true, nil
	}

	return resources.Label{}, false, errors.New("unable to parse response as error or valid response")
}

func (t Tracker) getLabelFromName(name string) (resources.Label, error) {
	labelsResponse, err := t.Caller.Get(LabelsEndpoint)
	if err != nil {
		return resources.Label{}, err
	}
	var labels resources.Labels
	if err := json.Unmarshal(labelsResponse, &labels); err != nil {
		fmt.Println("UnmarshalError")
		return resources.Label{}, err
	}

	for _, label := range labels {
		if label.Name == name {
			return label, nil
		}
	}

	return resources.Label{}, errors.New(fmt.Sprintf("no labels found with name: %s", name))
}

func (t Tracker) initializeStaleLabel() (bool, error) {
	return false, nil
}

func tagAsStale(story resources.Story) bool {
	return false
}
