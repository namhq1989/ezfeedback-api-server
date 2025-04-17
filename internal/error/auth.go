package apperrors

import "errors"

var Auth = struct {
	InvalidAuthToken    error
	InvalidRefreshToken error
	InvalidExpiry       error
	NotAllowed          error
}{
	InvalidAuthToken:    errors.New("auth_invalid_auth_token"),
	InvalidRefreshToken: errors.New("auth_invalid_refresh_token"),
	InvalidExpiry:       errors.New("auth_invalid_expiry"),
	NotAllowed:          errors.New("auth_not_allowed"),
}
