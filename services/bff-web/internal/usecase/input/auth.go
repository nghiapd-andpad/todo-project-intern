// Package input defines the data structures for input parameters used in use cases.
package input

type RegisterInput struct {
	Username string
	Password string
	Email    string
}

type LoginInput struct {
	Username string
	Password string
}
