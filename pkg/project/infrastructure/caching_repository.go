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

	domain                               string
	projectByIDCachingTime               time.Duration
	projectSettingByProjectIDCachingTime time.Duration
}

func NewCachingRepository(caching *caching.Caching, isEnvRelease bool) CachingRepository {
	if isEnvRelease {
		return CachingRepository{
			caching:                              caching,
			domain:                               "project",
			projectByIDCachingTime:               12 * time.Hour,
			projectSettingByProjectIDCachingTime: 12 * time.Hour,
		}
	} else {
		cachingTime := 1 * time.Minute

		return CachingRepository{
			caching:                              caching,
			domain:                               "project",
			projectByIDCachingTime:               cachingTime,
			projectSettingByProjectIDCachingTime: cachingTime,
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
