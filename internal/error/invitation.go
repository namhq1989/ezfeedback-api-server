package apperrors

import "errors"

var Invitation = struct {
	InvalidInvitationID error
	InvitationNotFound  error
	InvalidInviterID    error
	InvalidToken        error
}{
	InvalidInvitationID: errors.New("invitation_invalid_id"),
	InvitationNotFound:  errors.New("invitation_not_found"),
	InvalidInviterID:    errors.New("invitation_invalid_inviter_id"),
	InvalidToken:        errors.New("invitation_invalid_token"),
}
