package client

import (
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type UserClientImpl struct {
	UserRepo *crud.Repository[model.UserView, string]
}

func (u UserClientImpl) GetUserById(id string) model.UserView {
	record, err := u.UserRepo.FindById(id)
	if err != nil {
		return model.UserView{}
	}

	return *record
}

func (u UserClientImpl) GetUsersByIds(ids []string) []*model.UserView {
	users, err := u.UserRepo.FindByIds(ids)

	if err != nil {
		return nil
	}

	return users
}
