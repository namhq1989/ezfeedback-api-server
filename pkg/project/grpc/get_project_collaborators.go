package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectCollaboratorsHandler struct {
	service domain.Service
}

func NewGetProjectCollaboratorsHandler(service domain.Service) GetProjectCollaboratorsHandler {
	return GetProjectCollaboratorsHandler{
		service: service,
	}
}

func (h GetProjectCollaboratorsHandler) GetProjectCollaborators(ctx *appcontext.AppContext, req *projectpb.GetProjectCollaboratorsRequest) (*projectpb.GetProjectCollaboratorsResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project collaborators request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("find collaborators in db")
	collaborators, err := h.service.GetProjectCollaboratorsByProjectID(ctx, req.GetProjectId())
	if err != nil {
		ctx.Logger().Error("failed to find collaborators in db", err, appcontext.Fields{})
		return nil, err
	}

	var result = make([]*projectpb.ProjectCollaborator, 0)
	for _, collaborator := range collaborators {
		result = append(result, &projectpb.ProjectCollaborator{
			Id:     collaborator.ID,
			UserId: collaborator.UserID,
			Role:   collaborator.Role.String(),
		})
	}

	ctx.Logger().Text("done get project collaborators request")
	return &projectpb.GetProjectCollaboratorsResponse{
		Collaborators: result,
	}, nil
}
