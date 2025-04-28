package query

import (
	"sync"

	"github.com/goccy/go-json"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectByIDHandler struct {
	cachingRepository domain.CachingRepository
	service           domain.Service
}

func NewGetProjectByIDHandler(cachingRepository domain.CachingRepository, service domain.Service) GetProjectByIDHandler {
	return GetProjectByIDHandler{
		cachingRepository: cachingRepository,
		service:           service,
	}
}

// GetProjectByID godoc
// @tags     Project
// @summary  Get project by id
// @id       project-get-by-id
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload query    dto.GetProjectByIDRequest true "Query"
// @success  200     {object} dto.GetProjectByIDResponse
// @router   /api/project/{id} [get]
func (h GetProjectByIDHandler) GetProjectByID(ctx *appcontext.AppContext, performerID, projectID string, _ dto.GetProjectByIDRequest) (*dto.GetProjectByIDResponse, error) {
	ctx.Logger().Info("new get project by id request", appcontext.Fields{"performerID": performerID, "projectID": projectID})

	ctx.Logger().Text("find api caching data")
	apiCachingData, err := h.cachingRepository.GetApiGetProjectByID(ctx, projectID)
	if apiCachingData != nil {
		ctx.Logger().Text("found api caching data")
		var result = dto.GetProjectByIDResponse{}
		if err = json.Unmarshal([]byte(*apiCachingData), &result); err == nil {
			ctx.Logger().Text("done get project by id request")
			return &result, nil
		} else {
			ctx.Logger().Error("failed to unmarshal api caching data", err, appcontext.Fields{})
		}
	}

	ctx.Logger().Text("find project in db")
	project, err := h.service.GetProjectByID(ctx, projectID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	var (
		setting    *domain.ProjectSetting
		campaigns  = make([]domain.ProjectCampaign, 0)
		categories = make([]domain.ProjectCategory, 0)

		wg = sync.WaitGroup{}
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		ctx.Logger().Text("find project setting in db")
		setting, err = h.service.GetProjectSettingByProjectID(ctx, projectID)
		if err != nil {
			ctx.Logger().Error("failed to find project setting in db", err, appcontext.Fields{})
		}
		if setting == nil {
			ctx.Logger().ErrorText("project setting not found")
			setting = domain.DefaultProjectSetting()
		}
	}()

	go func() {
		defer wg.Done()
		ctx.Logger().Text("find project campaigns in db")
		campaigns, err = h.service.GetProjectCampaignsByProjectID(ctx, projectID)
		if err != nil {
			ctx.Logger().Error("failed to find project campaigns in db", err, appcontext.Fields{})
		}
	}()

	go func() {
		defer wg.Done()
		ctx.Logger().Text("find project categories in db")
		categories, err = h.service.GetProjectCategoriesByProjectID(ctx, projectID, domain.StatusUnknown)
		if err != nil {
			ctx.Logger().Error("failed to find project categories in db", err, appcontext.Fields{})
		}
	}()

	wg.Wait()

	ctx.Logger().Text("convert to response")
	result := dto.GetProjectByIDResponse{
		Project: dto.Project{}.FromDomain(*project, *setting, campaigns, categories),
	}

	ctx.Logger().Text("set api caching data")
	apiCachingDataByte, err := json.Marshal(result)
	if err != nil {
		ctx.Logger().Error("failed to marshal api caching data", err, appcontext.Fields{})
	} else {
		if err = h.cachingRepository.SetApiGetProjectByID(ctx, projectID, string(apiCachingDataByte)); err != nil {
			ctx.Logger().Error("failed to set api response to caching", err, appcontext.Fields{})
		}
	}

	ctx.Logger().Text("done get project by id request")
	return &result, nil
}
