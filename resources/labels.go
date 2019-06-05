package resources

type Labels []Label

type Label struct {
	Kind      string    `json:"kind,omitempty"`
	ID        int       `json:"id,omitempty"`
	ProjectID int       `json:"project_id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
