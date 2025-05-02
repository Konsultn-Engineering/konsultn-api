package model

type UserView struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

func (UserView) TableName() string {
	return "users"
}
