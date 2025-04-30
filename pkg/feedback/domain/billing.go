package domain

import "github.com/namhq1989/go-utilities/appcontext"

type BillingHub interface {
	CanAcceptFeedback(ctx *appcontext.AppContext, projectID string, totalCreatedFeedbacks int64) (bool, error)
}
