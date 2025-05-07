package domain

import "github.com/namhq1989/go-utilities/appcontext"

type IAMHub interface {
	GetUserByID(ctx *appcontext.AppContext, userID string) (*User, error)
}

type User struct {
	ID    string
	Name  string
	Email string
}
