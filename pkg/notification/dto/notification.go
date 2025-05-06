package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
)

type Notification struct {
	ID        string                    `json:"id"`
	Content   string                    `json:"content"`
	IsRead    bool                      `json:"isRead"`
	Action    NotificationAction        `json:"action"`
	CreatedAt *httprespond.TimeResponse `json:"createdAt"`
}

type NotificationAction struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

func (Notification) FromDomain(notification domain.Notification, content string, action domain.NotificationAction, params map[string]string) Notification {
	return Notification{
		ID:      notification.ID,
		Content: content,
		IsRead:  notification.IsRead,
		Action: NotificationAction{
			Action: action.String(),
			Params: params,
		},
		CreatedAt: httprespond.NewTimeResponse(notification.CreatedAt),
	}
}
