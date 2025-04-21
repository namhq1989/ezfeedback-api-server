package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetUserByID(ctx *appcontext.AppContext, id string) (*User, error)
	SetUserByID(ctx *appcontext.AppContext, id string, user User) error
	DeleteUserByID(ctx *appcontext.AppContext, id string) error
}
