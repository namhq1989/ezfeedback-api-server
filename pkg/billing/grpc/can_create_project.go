package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CanCreateProjectHandler struct{}

func NewCanCreateProjectHandler() CanCreateProjectHandler {
	return CanCreateProjectHandler{}
}

func (h CanCreateProjectHandler) CanCreateProject(ctx *appcontext.AppContext, req *billingpb.CanCreateProjectRequest) (*billingpb.CanCreateProjectResponse, error) {
	ctx.SetTraceID(req.TraceId)
	ctx.Logger().Info("new check can create project request", appcontext.Fields{"userId": req.UserId, "totalCreatedProjects": req.TotalCreatedProjects})

	ctx.Logger().Text("done check can create project request")
	return &billingpb.CanCreateProjectResponse{
		Allowed: true,
	}, nil
}
