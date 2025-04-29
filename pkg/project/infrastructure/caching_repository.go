package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CachingRepository struct {
	caching caching.Operations

	domain                                  string
	projectByIDCachingTime                  time.Duration
	projectSettingByProjectIDCachingTime    time.Duration
	projectCategoriesByProjectIDCachingTime time.Duration
	projectCampaignsByProjectIDCachingTime  time.Duration

	apiGetProjectsByUserIDCachingTime time.Duration
	apiGetProjectByIDCachingTime      time.Duration
}

func NewCachingRepository(caching *caching.Caching, isEnvRelease bool) CachingRepository {
	if isEnvRelease {
		return CachingRepository{
			caching:                                 caching,
			domain:                                  "project",
			projectByIDCachingTime:                  12 * time.Hour,
			projectSettingByProjectIDCachingTime:    12 * time.Hour,
			projectCategoriesByProjectIDCachingTime: 12 * time.Hour,
			projectCampaignsByProjectIDCachingTime:  12 * time.Hour,

			apiGetProjectsByUserIDCachingTime: 12 * time.Hour,
			apiGetProjectByIDCachingTime:      12 * time.Hour,
		}
	} else {
		cachingTime := 1 * time.Minute

		return CachingRepository{
			caching:                                 caching,
			domain:                                  "project",
			projectByIDCachingTime:                  cachingTime,
			projectSettingByProjectIDCachingTime:    cachingTime,
			projectCategoriesByProjectIDCachingTime: cachingTime,
			projectCampaignsByProjectIDCachingTime:  cachingTime,

			apiGetProjectsByUserIDCachingTime: cachingTime,
			apiGetProjectByIDCachingTime:      cachingTime,
		}
	}
}

//
// GET PROJECT BY ID
//

func (r CachingRepository) GetProjectByID(ctx *appcontext.AppContext, id string) (*domain.Project, error) {
	key := r.generateProjectByIDKey(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result *domain.Project
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetProjectByID(ctx *appcontext.AppContext, id string, project domain.Project) error {
	key := r.generateProjectByIDKey(id)
	r.caching.SetTTL(ctx, key, project, r.projectByIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteProjectByID(ctx *appcontext.AppContext, id string) error {
	key := r.generateProjectByIDKey(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateProjectByIDKey(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("project:%s", id))
}

//
// GET PROJECT SETTING BY PROJECT ID
//

func (r CachingRepository) GetProjectSettingByProjectID(ctx *appcontext.AppContext, id string) (*domain.ProjectSetting, error) {
	key := r.generateProjectSettingByProjectIDKey(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result *domain.ProjectSetting
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetProjectSettingByProjectID(ctx *appcontext.AppContext, id string, setting domain.ProjectSetting) error {
	key := r.generateProjectSettingByProjectIDKey(id)
	r.caching.SetTTL(ctx, key, setting, r.projectSettingByProjectIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteProjectSettingByProjectID(ctx *appcontext.AppContext, id string) error {
	key := r.generateProjectSettingByProjectIDKey(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateProjectSettingByProjectIDKey(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("project:%s:setting", id))
}

//
// GET PROJECT CATEGORIES BY PROJECT ID
//

func (r CachingRepository) GetProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string) ([]domain.ProjectCategory, error) {
	key := r.generateProjectCategoriesByProjectIDKey(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result []domain.ProjectCategory
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string, categories []domain.ProjectCategory) error {
	key := r.generateProjectCategoriesByProjectIDKey(id)
	r.caching.SetTTL(ctx, key, categories, r.projectCategoriesByProjectIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string) error {
	key := r.generateProjectCategoriesByProjectIDKey(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateProjectCategoriesByProjectIDKey(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("project:%s:categories", id))
}

//
// GET PROJECT CAMPAIGNS BY PROJECT ID
//

func (r CachingRepository) GetProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string) ([]domain.ProjectCampaign, error) {
	key := r.generateProjectCampaignsByProjectIDKey(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result []domain.ProjectCampaign
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string, campaigns []domain.ProjectCampaign) error {
	key := r.generateProjectCampaignsByProjectIDKey(id)
	r.caching.SetTTL(ctx, key, campaigns, r.projectCampaignsByProjectIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string) error {
	key := r.generateProjectCampaignsByProjectIDKey(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateProjectCampaignsByProjectIDKey(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("project:%s:campaigns", id))
}

//
// API GET PROJECTS BY USER ID
//

func (r CachingRepository) GetApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string) (*string, error) {
	key := r.generateApiGetProjectsByUserID(userID)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if dataStr == "" {
		return nil, nil
	}

	return &dataStr, nil
}

func (r CachingRepository) SetApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string, data string) error {
	key := r.generateApiGetProjectsByUserID(userID)
	r.caching.SetTTL(ctx, key, data, r.apiGetProjectsByUserIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string) error {
	key := r.generateApiGetProjectsByUserID(userID)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateApiGetProjectsByUserID(userID string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("api:getProjectsByUserId:%s", userID))
}

//
// API GET PROJECT BY ID
//

func (r CachingRepository) GetApiGetProjectByID(ctx *appcontext.AppContext, id string) (*string, error) {
	key := r.generateApiGetProjectByID(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if dataStr == "" {
		return nil, nil
	}

	return &dataStr, nil
}

func (r CachingRepository) SetApiGetProjectByID(ctx *appcontext.AppContext, id string, data string) error {
	key := r.generateApiGetProjectByID(id)
	r.caching.SetTTL(ctx, key, data, r.apiGetProjectByIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteApiGetProjectByID(ctx *appcontext.AppContext, id string) error {
	key := r.generateApiGetProjectByID(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateApiGetProjectByID(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("api:getProjectById:%s", id))
}
