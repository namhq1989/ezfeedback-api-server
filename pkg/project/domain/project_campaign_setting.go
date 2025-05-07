package domain

type ProjectCampaignSetting struct {
	WidgetPosition   string
	AllowAnonymous   bool
	EnableRating     bool
	MinRating        int32
	MaxRating        int32
	FollowUpQuestion string
}
