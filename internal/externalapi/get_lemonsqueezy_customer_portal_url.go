package externalapi

import (
	"fmt"

	"github.com/namhq1989/go-utilities/appcontext"
)

type GetLemonsqueezyCustomerPortalURLResult struct {
	URL string
}

type getLemonsqueezyCustomerDataResponse struct {
	Data getLemonsqueezyCustomerDataResponseData `json:"data"`
}

type getLemonsqueezyCustomerDataResponseData struct {
	ID         string                                            `json:"id"`
	Attributes getLemonsqueezyCustomerDataResponseDataAttributes `json:"attributes"`
}

type getLemonsqueezyCustomerDataResponseDataAttributes struct {
	URLs getLemonsqueezyCustomerDataResponseDataAttributesURL `json:"urLs"`
}

type getLemonsqueezyCustomerDataResponseDataAttributesURL struct {
	CustomerPortal string `json:"customer_portal"`
}

func (ea ExternalApi) GetLemonsqueezyCustomerPortalURL(ctx *appcontext.AppContext, customerID string) (*GetLemonsqueezyCustomerPortalURLResult, error) {
	var apiResult = &getLemonsqueezyCustomerDataResponse{}

	_, err := ea.lemonsqueezyClient.R().
		SetResult(&apiResult).
		Get(fmt.Sprintf("/v1/customers/%s", customerID))

	if err != nil {
		ctx.Logger().Error("[externalapi] error when get Lemonsqueezy customer data", err, appcontext.Fields{})
		return nil, err
	}

	return &GetLemonsqueezyCustomerPortalURLResult{
		URL: apiResult.Data.Attributes.URLs.CustomerPortal,
	}, nil
}
