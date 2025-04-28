package dto

type ProjectDTO struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Tasks []TaskDTO `json:"tasks,omitempty"`
}

type TaskDTO struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	AssigneeID  *string      `json:"assignee_id"`
	Assignee    *AssigneeDTO `json:"assignee"`
	ProjectID   string       `json:"project_id"`
}

type AssigneeDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
