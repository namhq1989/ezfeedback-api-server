package externalapi

import (
	"fmt"
	"strings"

	"github.com/namhq1989/go-utilities/appcontext"
)

type GetIpLocationDataResult struct {
	Country string
}

type getIpLocationDataApiResult struct {
	Country string `json:"country"`
}

func (ea ExternalApi) GetIpLocationData(ctx *appcontext.AppContext, ip string) (*GetIpLocationDataResult, error) {
	var apiResult = &getIpLocationDataApiResult{}

	_, err := ea.locationClient.R().
		SetQueryParams(map[string]string{
			"token": ea.ipInfoToken,
		}).
		SetResult(&apiResult).
		Get(fmt.Sprintf("/%s", ip))

	if err != nil {
		ctx.Logger().Error("[externalapi] error when get ip location data", err, appcontext.Fields{})
		return nil, err
	}

	if apiResult == nil || apiResult.Country == "" {
		return nil, nil
	}

	return &GetIpLocationDataResult{
		Country: strings.ToLower(apiResult.Country),
	}, nil
}
