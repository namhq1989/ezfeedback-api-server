package infrastructure

import (
	appjwt "github.com/namhq1989/ezfeedback-api-server/internal/jwt"
	"github.com/namhq1989/go-utilities/appcontext"
)

type JwtRepository struct {
	jwt appjwt.Operations
}

func NewJwtRepository(jwt appjwt.Operations) JwtRepository {
	return JwtRepository{
		jwt: jwt,
	}
}

func (r JwtRepository) GenerateAccessToken(ctx *appcontext.AppContext, userID string) (string, error) {
	return r.jwt.GenerateAccessToken(ctx, userID)
}
