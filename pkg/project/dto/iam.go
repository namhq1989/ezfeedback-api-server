package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (User) FromDomain(user domain.User) User {
	return User{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status.String(),
	}
}
