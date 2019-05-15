package resources

type Labels []Label

type Label struct {
	Kind      string `json:"kind,omitempty"`
	ID        int    `json:"id,omitempty"`
	ProjectID int    `json:"project_id,omitempty"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}
