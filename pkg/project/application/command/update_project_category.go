package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateProjectCategoryHandler struct {
	projectRepository         domain.ProjectRepository
	projectCategoryRepository domain.ProjectCategoryRepository
	cachingRepository         domain.CachingRepository
}

func NewUpdateProjectCategoryHandler(
	projectRepository domain.ProjectRepository,
	projectCategoryRepository domain.ProjectCategoryRepository,
	cachingRepository domain.CachingRepository,
) UpdateProjectCategoryHandler {
	return UpdateProjectCategoryHandler{
		projectRepository:         projectRepository,
		projectCategoryRepository: projectCategoryRepository,
		cachingRepository:         cachingRepository,
	}
}

// UpdateProjectCategory godoc
// @tags     Project
// @summary  Update project category
// @id       project-update-category
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    categoryId  	 path     string true "Category id"
// @param    payload body    dto.UpdateProjectCategoryRequest true "Body"
// @success  200     {object} dto.UpdateProjectCategoryResponse
// @router   /api/project/{projectId}/category/{categoryId} [put]
func (h UpdateProjectCategoryHandler) UpdateProjectCategory(ctx *appcontext.AppContext, performerID, projectID, categoryID string, req dto.UpdateProjectCategoryRequest) (*dto.UpdateProjectCategoryResponse, error) {
	ctx.Logger().Info("new update project category request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "categoryID": categoryID, "name": req.Name,
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

	ctx.Logger().Text("find project category in db")
	category, err := h.projectCategoryRepository.FindByID(ctx, categoryID)
	if err != nil {
		ctx.Logger().Error("failed to find project category in db", err, appcontext.Fields{})
		return nil, err
	}
	if category == nil {
		ctx.Logger().ErrorText("project category not found")
		return nil, apperrors.Project.InvalidCategory
	}
	if !category.IsBelongToProject(projectID) {
		ctx.Logger().ErrorText("project category not belong to project")
		return nil, apperrors.Project.InvalidCategory
	}

	ctx.Logger().Text("set project category data")
	if err = category.SetName(req.Name); err != nil {
		ctx.Logger().Error("failed to set project category data", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update project category in db")
	if err = h.projectCategoryRepository.Update(ctx, *category); err != nil {
		ctx.Logger().Error("failed to update project category in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCategoriesByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectsByUserID(ctx, performerID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done update project category request")
	return &dto.UpdateProjectCategoryResponse{}, nil
}
