package apperrors

import "errors"

var User = struct {
	InvalidUserID    error
	UserNotFound     error
	InvalidInvitedBy error
}{
	InvalidUserID:    errors.New("user_invalid_id"),
	UserNotFound:     errors.New("user_not_found"),
	InvalidInvitedBy: errors.New("user_invalid_invited_by"),
}
