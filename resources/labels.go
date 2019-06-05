package resources

import "time"

type Labels []Label

type Label struct {
	Kind      string    `json:"kind,omitempty"`
	ID        int       `json:"id,omitempty"`
	ProjectID int       `json:"project_id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
