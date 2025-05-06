package apperrors

import "errors"

var Notification = struct {
	NotificationNotFound                   error
	InvalidID                              error
	InvalidMetadata                        error
	InvalidType                            error
	UserProjectNotificationSettingNotFound error
}{
	NotificationNotFound:                   errors.New("notification_not_found"),
	InvalidID:                              errors.New("notification_invalid_id"),
	InvalidMetadata:                        errors.New("notification_invalid_metadata"),
	InvalidType:                            errors.New("notification_invalid_type"),
	UserProjectNotificationSettingNotFound: errors.New("notification_user_project_notification_setting_not_found"),
}
