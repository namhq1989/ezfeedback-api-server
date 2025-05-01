package externalapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Operations interface {
	GetIpLocationData(ctx *appcontext.AppContext, ip string) (*GetIpLocationDataResult, error)
	GenerateLemonsqueezySubscriptionCheckoutURL(ctx *appcontext.AppContext, userID, subscriptionID string) (*string, error)
	GetLemonsqueezySubscriptionInvoiceData(ctx *appcontext.AppContext, invoiceID string) (*GetLemonsqueezySubscriptionDataResult, error)
	GetLemonsqueezyCustomerPortalURL(ctx *appcontext.AppContext, customerID string) (*GetLemonsqueezyCustomerPortalURLResult, error)
}

type LemonsqueezyCfg struct {
	Token               string
	StoreID             string
	MonthlyVariantID    string
	MonthlyDiscountCode string
	YearlyVariantID     string
	YearlyDiscountCode  string
}

type ExternalApi struct {
	ipInfoToken     string
	lemonsqueezyCfg LemonsqueezyCfg

	locationClient     *resty.Client
	lemonsqueezyClient *resty.Client
}

const (
	locationApiEndpoint  = "https://ipinfo.io"
	lemonsqueezyEndpoint = "https://api.lemonsqueezy.com"
)

func NewExternalAPIClient(ipInfoToken string, lemonsqueezyCfg LemonsqueezyCfg) *ExternalApi {
	return &ExternalApi{
		ipInfoToken:     ipInfoToken,
		lemonsqueezyCfg: lemonsqueezyCfg,

		locationClient: resty.New().
			SetBaseURL(locationApiEndpoint).
			SetHeader("Accept", "application/json").
			SetTimeout(30 * time.Second).
			SetJSONMarshaler(json.Marshal).
			SetJSONUnmarshaler(json.Unmarshal).
			SetRetryAfter(func(_ *resty.Client, resp *resty.Response) (time.Duration, error) {
				return 1, fmt.Errorf("failed to send Location request at %s with status code %d", resp.Request.RawRequest.RequestURI, resp.StatusCode())
			}),
		lemonsqueezyClient: resty.New().
			SetBaseURL(lemonsqueezyEndpoint).
			SetHeader("Accept", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", lemonsqueezyCfg.Token)).
			SetTimeout(30 * time.Second).
			SetJSONMarshaler(json.Marshal).
			SetJSONUnmarshaler(json.Unmarshal).
			SetRetryAfter(func(_ *resty.Client, resp *resty.Response) (time.Duration, error) {
				return 1, fmt.Errorf("failed to send LemonSqueezy request at %s with status code %d", resp.Request.RawRequest.RequestURI, resp.StatusCode())
			}),
	}
}
