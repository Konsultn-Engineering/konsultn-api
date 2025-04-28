package dto

import (
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/mapper"
)

func mapTeamMembers(members []model.TeamMember) []TeamMemberDTO {
	result := make([]TeamMemberDTO, 0, len(members))
	for _, m := range members {
		result = append(result, TeamMemberDTO{
			ID:        m.User.ID,
			FirstName: m.User.FirstName,
			LastName:  m.User.LastName,
			Email:     m.User.Email,
		})
	}
	return result
}

func ToTeamDTO(team *model.Team) TeamDTO {
	dto, _ := mapper.Convert[TeamDTO](team)

	if team.Members != nil {
		mapped := mapTeamMembers(team.Members)
		dto.Members = &mapped
	}

	return dto
}
