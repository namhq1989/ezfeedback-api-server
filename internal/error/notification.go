package apperrors

import "errors"

var Notification = struct {
	InvalidMetadata                        error
	InvalidType                            error
	UserProjectNotificationSettingNotFound error
}{
	InvalidMetadata:                        errors.New("notification_invalid_metadata"),
	InvalidType:                            errors.New("notification_invalid_type"),
	UserProjectNotificationSettingNotFound: errors.New("notification_user_project_notification_setting_not_found"),
}
