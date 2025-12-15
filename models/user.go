package models

import (
	"errors"
	"net/mail"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Validate() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return err
	}
	if u.Name == "" {
		return errors.New("field 'name' cannot be empty")
	}
	if u.Email == "" {
		return errors.New("field 'email' cannot be empty")
	}
	return nil
}
