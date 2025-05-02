package apperrors

import "errors"

var Notification = struct {
	UserProjectNotificationSettingNotFound error
}{
	UserProjectNotificationSettingNotFound: errors.New("notification_user_project_notification_setting_not_found"),
}
