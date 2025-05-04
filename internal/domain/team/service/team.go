package service

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

func (s *TeamService) CreateTeam(createTeamRequest dto.CreateTeamRequest) (model.Team, error) {
	var createdTeam model.Team

	err := s.db.Transaction(func(tx *gorm.DB) error {
		createdTeam = model.Team{
			Name:      createTeamRequest.Name,
			Slug:      createTeamRequest.Slug,
			OwnerID:   s.actingUserId,
			UpdatedBy: s.actingUserId,
		}

		if err := tx.Create(&createdTeam).Error; err != nil {
			return err
		}

		// Auto-assign creator as team owner
		member := &model.TeamMember{
			TeamID: createdTeam.ID,
			UserID: s.actingUserId,
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

	owner := s.userClient.GetUserById(createdTeam.OwnerID)
	createdTeam.Owner = &owner
	_ = s.hydrateTeam(&createdTeam)

	return createdTeam, err
}

func (s *TeamService) GetTeamById(teamId string) (model.Team, error) {
	teamRecord, err := s.teamRepo.FindById(teamId)

	if err != nil {
		return model.Team{}, err
	}

	owner := s.userClient.GetUserById(teamRecord.OwnerID)
	teamRecord.Owner = &owner
	_ = s.hydrateTeam(&teamRecord)

	return teamRecord, nil
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

	team.UpdatedBy = userId

	team, err = s.teamRepo.Save(team)

	if err != nil {
		return model.Team{}, err
	}

	s.hydrateTeam(&team)

	return team, nil
}

func (s *TeamService) GetAllUserTeams(params crud.QueryParams) ([]model.Team, error) {
	userId := s.actingUserId

	allTeams, err := s.teamRepo.Query(crud.AdvancedQuery{
		QueryParams: params,
		Joins: []crud.JoinClause{
			{Table: "team_members", On: "team_members.team_id = teams.id", JoinType: "JOIN"},
		},
		Wheres: []crud.WhereClause{
			{Query: "team_members.user_id = ?", Args: []any{userId}},
		},
		Preload: []string{},
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching user teams: %w", err)
	}

	// 1. Collect unique owner IDs
	ownerIDSet := make(map[string]struct{})
	for _, team := range allTeams {
		ownerIDSet[team.OwnerID] = struct{}{}
	}

	ownerIDs := make([]string, 0, len(ownerIDSet))
	for id := range ownerIDSet {
		ownerIDs = append(ownerIDs, id)
	}

	// 2. Fetch all users in a single call
	owners := s.userClient.GetUsersByIds(ownerIDs)

	// 3. Map user ID to user
	ownerMap := make(map[string]model.UserView)
	for _, user := range owners {
		ownerMap[user.ID] = user
	}

	// 4. Assign Owner to each team
	for i := range allTeams {
		owner, exists := ownerMap[allTeams[i].OwnerID]
		if exists {
			allTeams[i].Owner = &owner
		}
	}

	return allTeams, nil
}

func (s *TeamService) CanUpdateOrDeleteTeam(teamId string, userId string) bool {
	return s.teamMemberRepo.IsTeamAdmin(teamId, userId)
}
