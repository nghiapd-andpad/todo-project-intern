// Package input defines the input structures for user-related use cases in the core-user service.
package input

type UserRegister struct {
	Username string
	Password string
	Email    string
}

type UserLogin struct {
	Username string
	Password string
}
