package template

import "strings"

var newFeedbackContent = map[string]string{
	"en": "You've received new feedback for your project {{title}}",
	"vi": "Bạn đã nhận được phản hồi mới cho dự án {{title}}",
}

func NewFeedbackContent(language string, projectTitle string) string {
	return strings.ReplaceAll(newFeedbackContent[language], "{{title}}", projectTitle)
}
