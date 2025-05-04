package dto

type TeamDTO struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Owner       *TeamMemberDTO   `json:"owner"`
	Members     *[]TeamMemberDTO `json:"members,omitempty"`
}

type TeamMemberDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	JoinedAt  string `json:"joined_at,omitempty"`
	Role      string `json:"role,omitempty"`
}
