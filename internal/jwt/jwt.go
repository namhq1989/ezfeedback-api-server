package appjwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Operations interface {
	RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc

	GenerateAccessToken(ctx *appcontext.AppContext, userID string) (string, error)
	ParseAccessToken(ctx *appcontext.AppContext, token string) (*Claims, error)
}

const (
	defaultAccessTokenTTL = time.Minute * 15 // 15 minutes
)

type JWT struct {
	accessTokenSecret []byte
	accessTokenTTL    time.Duration
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type Result struct {
	AccessToken       string
	AccessTokenExpiry time.Time
}

func Init(accessTokenSecret string, accessTokenTTL time.Duration) (*JWT, error) {
	if accessTokenTTL.Seconds() == 0 {
		accessTokenTTL = defaultAccessTokenTTL
	}

	return &JWT{
		accessTokenSecret: []byte(accessTokenSecret),
		accessTokenTTL:    accessTokenTTL,
	}, nil
}
