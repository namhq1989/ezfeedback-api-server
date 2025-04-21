package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (User) FromDomain(user domain.User) User {
	return User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
