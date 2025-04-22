package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ChangeProjectStatusHandler struct {
	projectRepository domain.ProjectRepository
	cachingRepository domain.CachingRepository
}

func NewChangeProjectStatusHandler(projectRepository domain.ProjectRepository, cachingRepository domain.CachingRepository) ChangeProjectStatusHandler {
	return ChangeProjectStatusHandler{
		projectRepository: projectRepository,
		cachingRepository: cachingRepository,
	}
}

// ChangeProjectStatus godoc
// @tags     Project
// @summary  Change project status
// @id       project-change-status
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload body    dto.ChangeProjectStatusRequest true "Body"
// @success  200     {object} dto.ChangeProjectStatusResponse
// @router   /api/project/{id}/status [patch]
func (h ChangeProjectStatusHandler) ChangeProjectStatus(ctx *appcontext.AppContext, performerID, projectID string, req dto.ChangeProjectStatusRequest) (*dto.ChangeProjectStatusResponse, error) {
	ctx.Logger().Info("new change project status request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "status": req.Status,
	})

	ctx.Logger().Text("find project in db")
	project, err := h.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}
	if !project.IsOwner(performerID) {
		ctx.Logger().ErrorText("user is not project owner")
		return nil, apperrors.Common.NotFound
	}

	if project.Status.IsEqual(req.Status) {
		ctx.Logger().Text("project status not changed, respond")
		return &dto.ChangeProjectStatusResponse{}, nil
	}

	ctx.Logger().Text("set project new status")
	if err = project.SetStatus(req.Status); err != nil {
		ctx.Logger().Error("failed to set project new status", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update project in db")
	if err = h.projectRepository.Update(ctx, *project); err != nil {
		ctx.Logger().Error("failed to update project in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("set project in caching")
	if err = h.cachingRepository.SetProjectByID(ctx, projectID, *project); err != nil {
		ctx.Logger().Error("failed to set project in caching", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectsByUserID(ctx, performerID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done change project status request")
	return &dto.ChangeProjectStatusResponse{}, nil
}
