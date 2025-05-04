package domain

import "time"

type Feedback struct {
	ID           string
	Project      Project
	AppUserID    *string
	Email        *string
	Content      string
	Rating       int32
	CampaignType ProjectCampaignType
	CreatedAt    time.Time
}
