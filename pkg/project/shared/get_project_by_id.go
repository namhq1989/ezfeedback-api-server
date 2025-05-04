package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectByID(ctx *appcontext.AppContext, projectID, userID string) (*domain.Project, error) {
	ctx.Logger().Info("[service] get project by id", appcontext.Fields{"projectID": projectID, "userID": userID})

	ctx.Logger().Text("find project in caching")
	project, err := s.cachingRepository.GetProjectByID(ctx, projectID)
	if project != nil {
		ctx.Logger().Text("project found in caching, return")
		return project, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project not found in caching, find in db")
	project, err = s.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}
	if userID != "" && !project.IsOwner(userID) {
		ctx.Logger().ErrorText("user is not project owner")
		return nil, apperrors.Common.NotFound
	}

	ctx.Logger().Text("set project in caching")
	if err = s.cachingRepository.SetProjectByID(ctx, projectID, *project); err != nil {
		ctx.Logger().Error("failed to set project in caching", err, appcontext.Fields{})
	}
	return project, nil
}
