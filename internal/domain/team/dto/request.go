package dto

import "konsultn-api/internal/domain/team/enum"

type CreateTeamRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateTeamRequest struct {
	Name    *string `json:"name,omitempty"`
	Slug    *string `json:"slug,omitempty"`
	OwnerId *string `json:"owner_id,omitempty"`
}

type AddMemberRequest struct {
	UserId  string    `json:"user_id" binding:"required"`
	Role    enum.Role `json:"role"`
	Message *string   `json:"message"`
}

type UpdateMemberRequest struct {
	Role enum.Role `json:"role" binding:"required"`
}
