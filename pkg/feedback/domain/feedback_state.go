package domain

type FeedbackState string

const (
	FeedbackStateUnknown    FeedbackState = ""
	FeedbackStateNew        FeedbackState = "new"
	FeedbackStateInReview   FeedbackState = "in_review"
	FeedbackStatePlanned    FeedbackState = "planned"
	FeedbackStateInProgress FeedbackState = "in_progress"
	FeedbackStateCompleted  FeedbackState = "completed"
	FeedbackStateDeclined   FeedbackState = "declined"
)

func (s FeedbackState) String() string {
	return string(s)
}

func (s FeedbackState) IsValid() bool {
	return s != FeedbackStateUnknown
}

func (s FeedbackState) IsEqual(state string) bool {
	dState := ToFeedbackState(state)
	return s == dState
}

func ToFeedbackState(s string) FeedbackState {
	switch s {
	case FeedbackStateNew.String():
		return FeedbackStateNew
	case FeedbackStateInReview.String():
		return FeedbackStateInReview
	case FeedbackStatePlanned.String():
		return FeedbackStatePlanned
	case FeedbackStateInProgress.String():
		return FeedbackStateInProgress
	case FeedbackStateCompleted.String():
		return FeedbackStateCompleted
	case FeedbackStateDeclined.String():
		return FeedbackStateDeclined
	default:
		return FeedbackStateUnknown
	}
}
