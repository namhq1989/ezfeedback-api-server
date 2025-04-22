package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type PingHandler struct{}

func NewPingHandler() PingHandler {
	return PingHandler{}
}

func (PingHandler) Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error) {
	ctx.Logger().Text("new ping request")
	return &dto.PingResponse{
		Success: true,
	}, nil
}
