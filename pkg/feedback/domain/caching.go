package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetIpLocationData(ctx *appcontext.AppContext, ip string) (*IpLocationData, error)
	SetIpLocationData(ctx *appcontext.AppContext, ip string, data IpLocationData) error

	GetFeedbackByID(ctx *appcontext.AppContext, id string) (*Feedback, error)
	SetFeedbackByID(ctx *appcontext.AppContext, id string, feedback Feedback) error
	DeleteFeedbackByID(ctx *appcontext.AppContext, id string) error
}
