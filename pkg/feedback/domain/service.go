package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type Service interface {
	GetIpLocationData(ctx *appcontext.AppContext, ip string) (*IpLocationData, error)
	GetFeedbackByID(ctx *appcontext.AppContext, feedbackID string) (*Feedback, error)
}
