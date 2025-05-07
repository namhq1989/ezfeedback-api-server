package mapping

type ProjectCampaignSetting struct {
	WidgetPosition   string `json:"widgetPosition"`
	AllowAnonymous   bool   `json:"allowAnonymous"`
	EnableRating     bool   `json:"enableRating"`
	MinRating        int32  `json:"minRating"`
	MaxRating        int32  `json:"maxRating"`
	FollowUpQuestion string `json:"followUpQuestion"`
}
