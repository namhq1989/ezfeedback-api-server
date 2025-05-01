package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

// TODO:
// - set domain caching repository
// - set domain external api repository

type CachingRepository struct {
	caching caching.Operations

	domain                    string
	ipLocationDataCachingTime time.Duration
}

func NewCachingRepository(caching *caching.Caching, isEnvRelease bool) CachingRepository {
	if isEnvRelease {
		return CachingRepository{
			caching:                   caching,
			domain:                    "feedback",
			ipLocationDataCachingTime: 12 * time.Hour,
		}
	} else {
		// cachingTime := 1 * time.Minute

		return CachingRepository{
			caching:                   caching,
			domain:                    "feedback",
			ipLocationDataCachingTime: 12 * time.Hour,
		}
	}
}

// IP LOCATION DATA

func (r CachingRepository) GetIpLocationData(ctx *appcontext.AppContext, ip string) (*domain.IpLocationData, error) {
	key := r.generateIpLocationDataKey(ip)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if dataStr == "" {
		return nil, nil
	}

	var result *domain.IpLocationData
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r CachingRepository) SetIpLocationData(ctx *appcontext.AppContext, ip string, data domain.IpLocationData) error {
	key := r.generateIpLocationDataKey(ip)
	r.caching.SetTTL(ctx, key, data, r.ipLocationDataCachingTime)
	return nil
}

func (r CachingRepository) generateIpLocationDataKey(ip string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("ip:%s:locationData", ip))
}
