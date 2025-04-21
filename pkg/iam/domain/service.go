package domain

import "github.com/namhq1989/go-utilities/appcontext"

type Service interface {
	GetUserByID(ctx *appcontext.AppContext, id string) (*User, error)
}
