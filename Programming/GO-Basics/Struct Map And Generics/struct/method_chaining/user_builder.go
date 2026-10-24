package main

import (
	"errors"
	"fmt"
)

type UserBuilder struct {
	Username string
	Email    string
	Err      error
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (u *UserBuilder) SetUserName(username string) *UserBuilder {
	if username == "" {
		u.Err = errors.New("username is empty")
	}
	u.Username = username
	return u
}

func (u *UserBuilder) SetEmail(email string) *UserBuilder {
	if email == "" {
		u.Err = errors.New("email is empty")
	}
	u.Email = email
	return u
}

func (u *UserBuilder) Save() error {
	if u.Err != nil {
		return u.Err
	}

	fmt.Println("User saved to DB")

	return nil
}

func Main() {

	err := NewUserBuilder().SetUserName("abc").SetEmail("ere").Save()

	err = NewUserBuilder().SetUserName("").SetEmail("dgd").Save()
	if err != nil {
		fmt.Println(err)
	}
}
