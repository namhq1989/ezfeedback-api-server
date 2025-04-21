package apperrors

import "github.com/pkg/errors"

var Billing = struct {
	ProjectLimitExceeded  error
	FeedbackLimitExceeded error
}{
	ProjectLimitExceeded:  errors.New("billing_project_limit_exceeded"),
	FeedbackLimitExceeded: errors.New("billing_feedback_limit_exceeded"),
}
