package output

type UserDTO struct {
	ID       string
	Username string
	Email    string
}

type UserRegister struct {
	User *UserDTO
}

type UserLogin struct {
	AccessToken string
	User        *UserDTO
}
