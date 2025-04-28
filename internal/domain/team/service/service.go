package service

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamService struct {
	teamRepo       *crud.Repository[model.Team]
	teamMemberRepo *crud.Repository[model.TeamMember]
	db             *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{
		teamRepo:       crud.NewRepository[model.Team](db),
		teamMemberRepo: crud.NewRepository[model.TeamMember](db),
		db:             db,
	}
}

func (s *TeamService) CreateTeam(userId string, createTeamRequest dto.CreateTeamRequest) (dto.TeamDTO, error) {
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
		return dto.TeamDTO{}, err
	}

	s.teamRepo.Preload(&createdTeam, []string{"Owner", "Members.User"}, "id", createdTeam.ID)

	return dto.ToTeamDTO(&createdTeam), err
}
