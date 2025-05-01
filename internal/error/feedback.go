package apperrors

import "errors"

var Feedback = struct {
	InvalidFeedbackID  error
	FeedbackNotFound   error
	InvalidContent     error
	InvalidRating      error
	InvalidState       error
	DailyLimitExceeded error
}{
	InvalidFeedbackID:  errors.New("feedback_invalid_id"),
	FeedbackNotFound:   errors.New("feedback_not_found"),
	InvalidContent:     errors.New("feedback_invalid_content"),
	InvalidRating:      errors.New("feedback_invalid_rating"),
	InvalidState:       errors.New("feedback_invalid_state"),
	DailyLimitExceeded: errors.New("feedback_daily_limit_exceeded"),
}
