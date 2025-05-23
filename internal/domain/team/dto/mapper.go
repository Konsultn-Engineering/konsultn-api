package dto

import (
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
	"konsultn-api/internal/shared/dto"
	"konsultn-api/internal/shared/mapper"
)

func mapTeamMembers(members []model.TeamMember) []TeamMember {
	result := make([]TeamMember, 0, len(members))
	for _, m := range members {
		result = append(result, TeamMember{
			ID:        m.User.ID,
			FirstName: m.User.FirstName,
			LastName:  m.User.LastName,
			Email:     m.User.Email,
			JoinedAt:  m.JoinedAt.String(),
			Role:      m.Role,
		})
	}
	return result
}

func ToTeamDTO(team model.Team) Team {
	dto, _ := mapper.Convert[Team](team)

	if team.Members != nil {
		mapped := mapTeamMembers(team.Members)
		dto.Members = &mapped
	}

	return dto
}

func ToTeamDTOPaginated(teams *crud.PaginatedResult[model.Team]) dto.PaginatedResultDTO[Team] {
	dataObject, err := dto.MapPaginatedResult(teams, func(team model.Team) Team {
		return ToTeamDTO(team)
	})

	if err != nil {
	}
	return dataObject
}
