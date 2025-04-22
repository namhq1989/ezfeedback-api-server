package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Commands interface {
		GenerateToken(ctx *appcontext.AppContext, req dto.GenerateTokenRequest) (*dto.GenerateTokenResponse, error)

		RequestVerificationCode(ctx *appcontext.AppContext, ip string, req dto.RequestVerificationCodeRequest) (*dto.RequestVerificationCodeResponse, error)
		VerifyVerificationCode(ctx *appcontext.AppContext, ip string, req dto.VerifyVerificationCodeRequest) (*dto.VerifyVerificationCodeResponse, error)

		UpdateMe(ctx *appcontext.AppContext, performerID string, req dto.UpdateMeRequest) (*dto.UpdateMeResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)

		GetMe(ctx *appcontext.AppContext, performerID string, _ dto.GetMeRequest) (*dto.GetMeResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.GenerateTokenHandler

		command.RequestVerificationCodeHandler
		command.VerifyVerificationCodeHandler

		command.UpdateMeHandler
	}
	queryHandlers struct {
		query.PingHandler

		query.GetMeHandler
	}
	Application struct {
		commandHandlers
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New(
	userRepository domain.UserRepository,
	verificationCodeRepository domain.VerificationCodeRepository,
	queueRepository domain.QueueRepository,
	jwtRepository domain.JwtRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			GenerateTokenHandler: command.NewGenerateTokenHandler(
				jwtRepository,
			),

			VerifyVerificationCodeHandler: command.NewVerifyVerificationCodeHandler(
				userRepository,
				verificationCodeRepository,
				jwtRepository,
			),
			RequestVerificationCodeHandler: command.NewRequestVerificationCodeHandler(
				verificationCodeRepository,
				queueRepository,
			),

			UpdateMeHandler: command.NewUpdateMeHandler(
				userRepository,
				cachingRepository,
				service,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),

			GetMeHandler: query.NewGetMeHandler(service),
		},
	}
}
