package dto

type Team struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	Owner       *TeamMember   `json:"owner"`
	Members     *[]TeamMember `json:"members,omitempty"`
}

type TeamMember struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	JoinedAt  string `json:"joined_at,omitempty"`
	Role      string `json:"role,omitempty"`
}

type TeamSummary struct {
	*Team
	Members     int `json:"-"`
	MemberCount int `json:"member_count"`
}
