package dto

type CountFeedbacksRequest struct {
	ProjectID    string `query:"projectId"`
	CategoryID   string `query:"categoryId"`
	CampaignType string `query:"campaignType"`
	Keyword      string `query:"keyword"`
	State        string `query:"state"`
	Rating       int32  `query:"rating"`
}

type CountFeedbacksResponse struct {
	Total int64 `json:"total"`
}
