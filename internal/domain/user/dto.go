package user

type CreateUserRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UpdateUserRequest struct {
}
