package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CachingRepository struct {
	caching caching.Operations

	domain              string
	userByIDCachingTime time.Duration
}

func NewCachingRepository(caching *caching.Caching, isEnvRelease bool) CachingRepository {
	if isEnvRelease {
		return CachingRepository{
			caching:             caching,
			domain:              "iam",
			userByIDCachingTime: 12 * time.Hour,
		}
	} else {
		cachingTime := 1 * time.Minute

		return CachingRepository{
			caching:             caching,
			domain:              "iam",
			userByIDCachingTime: cachingTime,
		}
	}
}

//
// GET USER BY ID
//

func (r CachingRepository) GetUserByID(ctx *appcontext.AppContext, id string) (*domain.User, error) {
	key := r.generateUserByIDKey(id)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result *domain.User
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetUserByID(ctx *appcontext.AppContext, id string, user domain.User) error {
	key := r.generateUserByIDKey(id)
	r.caching.SetTTL(ctx, key, user, r.userByIDCachingTime)
	return nil
}

func (r CachingRepository) DeleteUserByID(ctx *appcontext.AppContext, id string) error {
	key := r.generateUserByIDKey(id)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateUserByIDKey(id string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("user:%s", id))
}
