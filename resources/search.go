package resources

import "time"

type Story struct {
	Kind          string        `json:"kind"`
	ID            int           `json:"id"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Estimate      int           `json:"estimate,omitempty"`
	StoryType     string        `json:"story_type"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	CurrentState  string        `json:"current_state"`
	RequestedByID int           `json:"requested_by_id"`
	URL           string        `json:"url"`
	ProjectID     int           `json:"project_id"`
	OwnerIds      []interface{} `json:"owner_ids"`
	Labels        []Label `json:"labels"`
	OwnedByID int `json:"owned_by_id,omitempty"`
}

type Search struct {
	Stories struct {
		Stories              []Story `json:"stories"`
		TotalPoints          int     `json:"total_points"`
		TotalPointsCompleted int     `json:"total_points_completed"`
		TotalHits            int     `json:"total_hits"`
		TotalHitsWithDone    int     `json:"total_hits_with_done"`
	} `json:"stories"`
	Epics struct {
		Epics     []interface{} `json:"epics"`
		TotalHits int           `json:"total_hits"`
	} `json:"epics"`
	Query string `json:"query"`
}
