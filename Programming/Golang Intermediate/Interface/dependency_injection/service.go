package main

import "fmt"

type UserService struct {
	logger Logger
}

func NewUserService(l Logger) *UserService {
	return &UserService{
		logger: l,
	}
}

func (us *UserService) CreateUser(user string) {
	us.logger.Log("Creating user: " + user)
	fmt.Println("User created")
}