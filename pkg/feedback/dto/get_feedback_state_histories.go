package dto

type GetFeedbackStateHistoriesRequest struct {
	Page int64 `query:"page"`
}

type GetFeedbackStateHistoriesResponse struct {
	Histories []FeedbackStateHistory `json:"histories"`
	Limit     int64                  `json:"limit"`
}
