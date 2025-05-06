package dto

type GetFeedbacksRequest struct {
	ProjectID    string `query:"projectId"`
	CategoryID   string `query:"categoryId"`
	CampaignType string `query:"campaignType"`
	Keyword      string `query:"keyword"`
	State        string `query:"state"`
	Rating       int32  `query:"rating"`
	Page         int64  `query:"page"`
}

type GetFeedbacksResponse struct {
	Feedbacks []Feedback `json:"feedbacks"`
	Limit     int64      `json:"limit"`
}
