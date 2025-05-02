package dto

import "konsultn-api/internal/domain/team/enum"

type CreateTeamRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateTeamRequest struct {
	Name          *string            `json:"name,omitempty"`
	Slug          *string            `json:"slug,omitempty"`
	OwnerId       *string            `json:"owner_id,omitempty"`
	AddMembers    []AddMemberRequest `json:"add_members,omitempty"`
	RemoveMembers []string           `json:"remove_members,omitempty"`
}

type AddMemberRequest struct {
	UserId  string    `json:"user_id" binding:"required"`
	Role    enum.Role `json:"role"`
	Message *string   `json:"message"`
}
