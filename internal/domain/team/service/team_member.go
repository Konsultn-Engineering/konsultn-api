package service

import (
	"fmt"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/shared/crud/types"
)

func (s *TeamService) UpdateTeamMember(teamId string, memberId string, updateMemberRequest dto.UpdateMemberRequest) error {
	record, err := s.teamMemberRepo.FindWhere(map[string]interface{}{
		"team_id": teamId,
		"user_id": memberId,
	})

	if err != nil {
		return fmt.Errorf("unable to find team member %s", err)
	}

	if record[0].Role == updateMemberRequest.Role.String() {
		return nil
	}

	record[0].Role = updateMemberRequest.Role.String()
	record[0].UpdatedBy = s.actingUserId

	_, err = s.teamMemberRepo.Save(record[0])

	if err != nil {
		return fmt.Errorf("here")
	}

	return nil
}

func (s *TeamService) RemoveTeamMember(teamId string, memberId string) error {
	findWhere := map[string]interface{}{
		"team_id": teamId,
		"user_id": memberId,
	}

	teamMember, err := s.teamMemberRepo.FindWhere(findWhere)

	fmt.Println(teamMember)

	if err != nil {
		return nil
	}

	err = s.teamMemberRepo.SoftDeleteWithUpdate(teamMember[0], types.UpdateMap{
		"updated_by": s.actingUserId,
	})

	if err != nil {
		return err
	}

	return nil
}
