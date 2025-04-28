package dto

type TeamDTO struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Slug    string           `json:"slug"`
	Owner   *TeamMemberDTO   `json:"owner"`
	Members *[]TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreateTeamRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
