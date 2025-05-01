package domain

import "github.com/namhq1989/go-utilities/appcontext"

type BillingHub interface {
	CanCreateProject(ctx *appcontext.AppContext, userID string, totalCreatedProjects int64) (bool, error)
}
