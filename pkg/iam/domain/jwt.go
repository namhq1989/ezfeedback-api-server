package domain

import "github.com/namhq1989/go-utilities/appcontext"

type JwtRepository interface {
	GenerateAccessToken(ctx *appcontext.AppContext, userID string) (string, error)
}
