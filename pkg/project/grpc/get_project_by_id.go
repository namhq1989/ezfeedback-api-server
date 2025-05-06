package grpc

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectByIDHandler struct {
	service domain.Service
}

func NewGetProjectByIDHandler(service domain.Service) GetProjectByIDHandler {
	return GetProjectByIDHandler{
		service: service,
	}
}

func (h GetProjectByIDHandler) GetProjectById(ctx *appcontext.AppContext, req *projectpb.GetProjectByIdRequest) (*projectpb.GetProjectByIdResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project by id request", appcontext.Fields{"projectID": req.GetProjectId(), "userID": req.GetUserId()})

	ctx.Logger().Text("find project in db")
	project, err := h.service.GetProjectByID(ctx, req.GetProjectId(), req.GetUserId())
	if err != nil {
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("done get project by id request")
	return &projectpb.GetProjectByIdResponse{
		Project: &projectpb.Project{
			Id:      project.ID,
			UserId:  project.UserID,
			Title:   project.Title,
			Status:  project.Status.String(),
			Setting: nil,
		},
	}, nil
}
