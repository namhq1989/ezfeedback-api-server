package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (User) FromDomain(user domain.User) User {
	return User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
