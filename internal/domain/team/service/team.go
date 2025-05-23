package service

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
	"maps"
	"slices"
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

	team.UpdatedBy = userId

	team, err = s.teamRepo.Save(team)

	if err != nil {
		return model.Team{}, err
	}

	err = s.hydrateTeam(team)

	if err != nil {
		return model.Team{}, err
	}

	return *team, nil
}

func (s *TeamService) GetAllUserTeams(params crud.QueryParams) (*crud.PaginatedResult[model.Team], error) {
	userId := s.actingUserId

	qb := s.teamRepo.Members().
		Select("teams.id", "teams.name", []string{"COUNT(team_members.user_id)", "member_count"}).
		Join("team_members", "tm_user").OnGroup(
		func(jb crud.JoinBuilder) {
			jb.And("id", "=", "team_id")
			jb.And("tm_user.user_id", "=", jb.Raw(userId))
		}).
		GroupBy("teams.id").
		WithPageParams(params)

	allTeams, err := qb.Paginate()

	if err != nil {
		return nil, fmt.Errorf("error fetching user teams: %w", err)
	}

	// 1. Collect unique owner IDs
	ownerIDSet := make(map[string]struct{})
	for _, team := range allTeams.Result {
		ownerIDSet[team.OwnerID] = struct{}{}
	}

	ownerIDs := maps.Keys(ownerIDSet) // Requires Go 1.21+

	// 2. Fetch and map owners
	owners := s.userClient.GetUsersByIds(slices.Sorted(ownerIDs))
	ownerMap := make(map[string]*model.UserView, len(owners))
	for _, user := range owners {
		ownerMap[user.ID] = user
	}

	// 3. Assign owners to teams
	for _, team := range allTeams.Result {
		if owner, ok := ownerMap[team.OwnerID]; ok {
			team.Owner = owner
		}
	}

	return (*crud.PaginatedResult[model.Team])(allTeams), nil
}

func (s *TeamService) Testing(params crud.QueryParams) any {
	res, _ := s.teamRepo.Members().
		Select("teams.*", []string{"COUNT(team_members.user_id)", "member_count"}).
		Join("team_members", "tm_user").OnGroup(
		func(jb crud.JoinBuilder) {
			jb.On("id", "=", "team_id")
			jb.And("tm_user.user_id", "IN", jb.RawSQL("(SELECT id from users where id = ?)", s.actingUserId))
		}).
		GroupBy("teams.id").
		WithPageParams(params).PaginateMap()

	allTeams, _ := crud.ConvertPaginated[model.TeamSummaryView](res)

	ownerIDSet := make(map[string]struct{})
	for _, team := range allTeams.Result {
		ownerIDSet[team.OwnerID] = struct{}{}
	}

	ownerIDs := maps.Keys(ownerIDSet) // Requires Go 1.21+

	// 2. Fetch and map owners
	owners := s.userClient.GetUsersByIds(slices.Sorted(ownerIDs))
	ownerMap := make(map[string]*model.UserView, len(owners))
	for _, user := range owners {
		ownerMap[user.ID] = user
	}

	// 3. Assign owners to teams
	for _, team := range allTeams.Result {
		if owner, ok := ownerMap[team.OwnerID]; ok {
			team.Owner = owner
		}
	}

	return allTeams
}

func (s *TeamService) CanUpdateOrDeleteTeam(teamId string, userId string) bool {
	return s.teamMemberRepo.IsTeamAdmin(teamId, userId)
}
