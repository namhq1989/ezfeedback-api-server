package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/externalapi"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ExternalAPIRepository struct {
	ea externalapi.Operations
}

func NewExternalAPIRepository(ea externalapi.Operations) ExternalAPIRepository {
	return ExternalAPIRepository{
		ea: ea,
	}
}

func (r ExternalAPIRepository) GetIpLocationData(ctx *appcontext.AppContext, ip string) (*domain.IpLocationData, error) {
	result, err := r.ea.GetIpLocationData(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &domain.IpLocationData{
		Country: result.Country,
	}, nil
}
