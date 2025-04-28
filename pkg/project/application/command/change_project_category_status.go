package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ChangeProjectCategoryStatusHandler struct {
	projectCategoryRepository domain.ProjectCategoryRepository
	cachingRepository         domain.CachingRepository
	service                   domain.Service
}

func NewChangeProjectCategoryStatusHandler(
	projectCategoryRepository domain.ProjectCategoryRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) ChangeProjectCategoryStatusHandler {
	return ChangeProjectCategoryStatusHandler{
		projectCategoryRepository: projectCategoryRepository,
		cachingRepository:         cachingRepository,
		service:                   service,
	}
}

// ChangeProjectCategoryStatus godoc
// @tags     Project
// @summary  Change project category status
// @id       project-change-category-status
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    categoryId  	 path     string true "Category id"
// @param    payload body    dto.UpdateProjectCategoryRequest true "Body"
// @success  200     {object} dto.UpdateProjectCategoryResponse
// @router   /api/project/{projectId}/category/{categoryId}/status [patch]
func (h ChangeProjectCategoryStatusHandler) ChangeProjectCategoryStatus(ctx *appcontext.AppContext, performerID, projectID, categoryID string, req dto.ChangeProjectCategoryStatusRequest) (*dto.ChangeProjectCategoryStatusResponse, error) {
	ctx.Logger().Info("new change project category status request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "categoryID": categoryID, "status": req.Status,
	})

	ctx.Logger().Text("find project category in db")
	category, err := h.service.GetProjectCategory(ctx, projectID, categoryID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find project category in db", err, appcontext.Fields{})
		return nil, err
	}
	if category == nil {
		ctx.Logger().ErrorText("project category not found")
		return nil, apperrors.Project.InvalidCategory
	}

	if category.Status.IsEqual(req.Status) {
		ctx.Logger().Text("project category status not changed, respond")
		return &dto.ChangeProjectCategoryStatusResponse{}, nil
	}

	ctx.Logger().Text("set project category new status")
	if err = category.SetStatus(req.Status); err != nil {
		ctx.Logger().Error("failed to set project category new status", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update project category in db")
	if err = h.projectCategoryRepository.Update(ctx, *category); err != nil {
		ctx.Logger().Error("failed to update project category in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCategoriesByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to set project in caching", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectsByUserID(ctx, performerID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done change project category status request")
	return &dto.ChangeProjectCategoryStatusResponse{}, nil
}
