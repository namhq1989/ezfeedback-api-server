package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetIpLocationData(ctx *appcontext.AppContext, ip string) (*domain.IpLocationData, error) {
	ctx.Logger().Info("[service] get ip location data", appcontext.Fields{"ip": ip})

	ctx.Logger().Text("find in caching")
	data, err := s.cachingRepository.GetIpLocationData(ctx, ip)
	if data != nil {
		ctx.Logger().Text("found in caching, respond")
		return data, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("not found in caching, call api")
	data, err = s.externalAPIRepository.GetIpLocationData(ctx, ip)
	if err != nil {
		ctx.Logger().Error("failed to call api", err, appcontext.Fields{})
		return nil, err
	}
	if data == nil {
		ctx.Logger().Text("ip location not found")
		return nil, apperrors.Common.BadRequest
	}

	ctx.Logger().Text("persist in caching")
	if err = s.cachingRepository.SetIpLocationData(ctx, ip, *data); err != nil {
		ctx.Logger().Error("failed to persist in caching", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done get ip location data")
	return data, nil
}
