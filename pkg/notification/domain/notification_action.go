package domain

type NotificationAction string

const (
	NotificationActionUnknown                 NotificationAction = ""
	NotificationActionOpenProjectFeedbackPage NotificationAction = "open_project_feedback_page"
)

func (n NotificationAction) String() string {
	return string(n)
}

func GenerateActionDataFromNotificationType(notificationType NotificationType, metadata NotificationMetadata) (NotificationAction, map[string]string) {
	var (
		action NotificationAction
		params = map[string]string{}
	)

	if notificationType.IsNewFeedback() {
		action = NotificationActionOpenProjectFeedbackPage
		params = map[string]string{
			"projectId": metadata.ProjectID,
		}
	}

	return action, params
}
