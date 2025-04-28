package query

import (
	"github.com/goccy/go-json"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectsHandler struct {
	projectRepository domain.ProjectRepository
	cachingRepository domain.CachingRepository
	service           domain.Service
}

func NewGetProjectsHandler(projectRepository domain.ProjectRepository, cachingRepository domain.CachingRepository, service domain.Service) GetProjectsHandler {
	return GetProjectsHandler{
		projectRepository: projectRepository,
		cachingRepository: cachingRepository,
		service:           service,
	}
}

// GetProjects godoc
// @tags     Project
// @summary  Get projects
// @id       project-get-list
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.GetProjectsRequest true "Query"
// @success  200     {object} dto.GetProjectsResponse
// @router   /api/project [get]
func (h GetProjectsHandler) GetProjects(ctx *appcontext.AppContext, performerID string, _ dto.GetProjectsRequest) (*dto.GetProjectsResponse, error) {
	ctx.Logger().Info("new get projects request", appcontext.Fields{"performerID": performerID})

	ctx.Logger().Text("find api caching data")
	apiCachingData, err := h.cachingRepository.GetApiGetProjectsByUserID(ctx, performerID)
	if apiCachingData != nil {
		ctx.Logger().Text("found api caching data")
		var result = dto.GetProjectsResponse{}
		if err = json.Unmarshal([]byte(*apiCachingData), &result); err == nil {
			ctx.Logger().Text("done get projects request")
			return &result, nil
		} else {
			ctx.Logger().Error("failed to unmarshal api caching data", err, appcontext.Fields{})
		}
	}

	ctx.Logger().Text("find projects in db")
	projects, err := h.projectRepository.FindByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find projects in db", err, appcontext.Fields{})
		return nil, err
	}
	if projects == nil || len(projects) == 0 {
		ctx.Logger().Text("this user has no projects, respond")
		return &dto.GetProjectsResponse{Projects: make([]dto.ProjectBrief, 0)}, nil
	}

	ctx.Logger().Text("convert to dto")
	result := h.convertToDto(ctx, projects)

	ctx.Logger().Text("set api caching data")
	apiCachingDataByte, err := json.Marshal(result)
	if err != nil {
		ctx.Logger().Error("failed to marshal api caching data", err, appcontext.Fields{})
	} else {
		if err = h.cachingRepository.SetApiGetProjectsByUserID(ctx, performerID, string(apiCachingDataByte)); err != nil {
			ctx.Logger().Error("failed to set api response to caching", err, appcontext.Fields{})
		}
	}

	ctx.Logger().Text("done get projects request")
	return &result, nil
}

func (h GetProjectsHandler) convertToDto(ctx *appcontext.AppContext, projects []domain.Project) dto.GetProjectsResponse {
	var result = make([]dto.ProjectBrief, 0)

	for _, project := range projects {
		result = append(result, dto.ProjectBrief{}.FromDomain(project))
	}

	return dto.GetProjectsResponse{Projects: result}
}
