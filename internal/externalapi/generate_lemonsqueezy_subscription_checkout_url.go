package externalapi

import (
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
)

// https://docs.lemonsqueezy.com/api/checkouts/create-checkout

type generateLemonsqueezyCheckoutURLResponse struct {
	Data generateLemonsqueezyCheckoutURLResponseData `json:"data"`
}

type generateLemonsqueezyCheckoutURLResponseData struct {
	Attributes generateLemonsqueezyCheckoutURLResponseDataAttributes `json:"attributes"`
}

type generateLemonsqueezyCheckoutURLResponseDataAttributes struct {
	URL string `json:"url"`
}

func (ea ExternalApi) GenerateLemonsqueezySubscriptionCheckoutURL(ctx *appcontext.AppContext, userID, subscriptionID string) (*string, error) {
	var (
		variantID    = ea.lemonsqueezyCfg.MonthlyVariantID
		discountCode = ea.lemonsqueezyCfg.MonthlyDiscountCode

		apiResult = &generateLemonsqueezyCheckoutURLResponse{}
	)

	if subscriptionID == "year" {
		variantID = ea.lemonsqueezyCfg.YearlyVariantID
		discountCode = ea.lemonsqueezyCfg.YearlyDiscountCode
	}

	var payload = map[string]interface{}{
		"data": map[string]interface{}{
			"type": "checkouts",
			"attributes": map[string]interface{}{
				"checkout_data": map[string]interface{}{
					"custom": map[string]string{
						"user_id": userID,
					},
				},
				"checkout_options": map[string]interface{}{
					"embed":    true,
					"discount": true,
				},
				"expires_at": time.Now().Add(time.Hour * 1).Format(time.RFC3339),
				"preview":    true,
			},
			"relationships": map[string]interface{}{
				"store": map[string]interface{}{
					"data": map[string]interface{}{
						"type": "stores",
						"id":   ea.lemonsqueezyCfg.StoreID,
					},
				},
				"variant": map[string]interface{}{
					"data": map[string]interface{}{
						"type": "variants",
						"id":   variantID,
					},
				},
			},
		},
	}

	if discountCode != "" {
		attributes := payload["data"].(map[string]interface{})["attributes"].(map[string]interface{})
		checkoutData := attributes["checkout_data"].(map[string]interface{})
		checkoutData["discount_code"] = discountCode
	}

	_, err := ea.lemonsqueezyClient.R().
		SetBody(payload).
		SetResult(&apiResult).
		Post("/v1/checkouts")

	if err != nil {
		ctx.Logger().Error("[externalapi] error when create Lemonsqueezy checkout url", err, appcontext.Fields{})
		return nil, err
	}

	// ctx.Logger().Print("payload", payload)
	// ctx.Logger().Print("apiresult", apiResult)
	// ctx.Logger().Print("resp", string(resp.Body()))

	return &apiResult.Data.Attributes.URL, nil
}
