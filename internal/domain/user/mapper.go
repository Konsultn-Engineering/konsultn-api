package user

// ToUserModel
//@TODO: Create Generic mapper
///*

func ToUserModel(dto CreateUserRequest) User {
	return User{
		Email:    dto.Email,
		Password: dto.Password,
	}
	//var user *User
	//shared.Mapper(dto, &user)
	//
	//return user
}
