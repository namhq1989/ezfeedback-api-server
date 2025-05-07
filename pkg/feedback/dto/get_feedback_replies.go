package dto

type GetFeedbackRepliesRequest struct {
	Page int64 `query:"page"`
}

type GetFeedbackRepliesResponse struct {
	Replies []FeedbackReply `json:"replies"`
	Limit   int64           `json:"limit"`
}
