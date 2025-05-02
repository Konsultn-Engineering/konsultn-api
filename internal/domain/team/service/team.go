package service

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/model"
)

func (s *TeamService) CreateTeam(userId string, createTeamRequest dto.CreateTeamRequest) (model.Team, error) {
	var createdTeam model.Team

	err := s.db.Transaction(func(tx *gorm.DB) error {
		createdTeam = model.Team{
			Name:    createTeamRequest.Name,
			Slug:    createTeamRequest.Slug,
			OwnerID: userId,
		}

		if err := tx.Create(&createdTeam).Error; err != nil {
			return err
		}

		// Auto-assign creator as team owner
		member := &model.TeamMember{
			TeamID: createdTeam.ID,
			UserID: userId,
			Role:   "owner",
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return model.Team{}, err
	}

	createdTeam.Owner = s.userClient.GetUserById(createdTeam.OwnerID)
	_ = s.hydrateTeam(&createdTeam)

	return createdTeam, err
}

func (s *TeamService) GetTeamById(teamId string) (model.Team, error) {
	teamRecord, err := s.teamRepo.FindById(teamId)

	if err != nil {
		return model.Team{}, err
	}

	teamRecord.Owner = s.userClient.GetUserById(teamRecord.OwnerID)
	_ = s.hydrateTeam(teamRecord)

	return *teamRecord, nil
}

func (s *TeamService) UpdateTeamById(userId string, teamId string, updateTeamRequest dto.UpdateTeamRequest) (model.Team, error) {
	team, err := s.teamRepo.FindById(teamId)

	if err != nil {
		return model.Team{}, err
	}

	if updateTeamRequest.Name != nil {
		team.Name = *updateTeamRequest.Name
	}

	if updateTeamRequest.Slug != nil {
		team.Slug = *updateTeamRequest.Slug
	}

	if updateTeamRequest.OwnerId != nil {
		if userId != team.OwnerID {
			return model.Team{}, fmt.Errorf("only team owner can transfer ownership")
		}
		team.OwnerID = *updateTeamRequest.OwnerId
	}

	err = s.syncTeamMembers(teamId, updateTeamRequest.AddMembers, updateTeamRequest.RemoveMembers)
	if err != nil {
		return model.Team{}, err
	}

	team, err = s.teamRepo.Save(team)

	if err != nil {
		return model.Team{}, err
	}

	s.hydrateTeam(team)

	return *team, nil
}

func (s *TeamService) CanUpdateOrDeleteTeam(teamId string, userId string) bool {
	return s.teamMemberRepo.IsTeamAdmin(teamId, userId)
}
