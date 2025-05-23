package client

import "konsultn-api/internal/domain/user"

type UserClient interface {
	GetUserById(id string) user.User
	GetUsersByTeamId(teamId string) []*user.User
}
