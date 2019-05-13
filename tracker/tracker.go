package tracker

import (
	"encoding/json"
	"github.com/dwillist/stale_issues/resources"
)

type Tracker struct {
	Caller Caller
}

type Caller interface {
	Call(endpoint string) (string, error)
}

func NewTracker(caller Caller) Tracker {
	return Tracker{
		Caller: caller,
	}
}

func (t Tracker) FilterIssues() (resources.Search, error){
	var result resources.Search
	const endpoint = "https://www.pivotaltracker.com/services/v5/projects/1042066/search?query=label%3Agithub-issue%20AND%20-state%3Aaccepted%20-state%3Afinished%20-state%3Adelivered"

	response, err := t.Caller.Call(endpoint)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return result, err
	}

	return result, nil
}