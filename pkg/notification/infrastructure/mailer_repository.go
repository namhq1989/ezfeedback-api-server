package infrastructure

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/namhq1989/ezfeedback-api-server/internal/mailer"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

//go:embed template/new_feedback_email.html
var newFeedbackEmailTemplate []byte

//go:embed template/feedback_item.html
var feedbackItemTemplate []byte

type MailerRepository struct {
	mailer                   mailer.Operations
	newFeedbackEmailSubject  string
	newFeedbackEmailTemplate string
	feedbackItemTemplate     string
}

func NewMailerRepository(mailer mailer.Operations) MailerRepository {
	var m = MailerRepository{
		mailer: mailer,

		newFeedbackEmailSubject:  "[EasyFeedback] New feedback received",
		newFeedbackEmailTemplate: string(newFeedbackEmailTemplate),
		feedbackItemTemplate:     string(feedbackItemTemplate),
	}

	return m
}

func (r MailerRepository) SendNewFeedbackEmail(ctx *appcontext.AppContext, toEmail string, feedbacks []domain.Feedback) error {
	totalFeedbacks := len(feedbacks)
	if totalFeedbacks == 0 {
		ctx.Logger().ErrorText("no feedback found when sending new feedback email")
		return nil
	}

	var (
		subject = r.newFeedbackEmailSubject
		content = r.newFeedbackEmailTemplate
	)

	feedbackLabel := "piece of feedback"
	if totalFeedbacks > 1 {
		feedbackLabel = "pieces of feedback"
	}

	var feedbackItemsContent string

	maxToShow := 2
	if totalFeedbacks < maxToShow {
		maxToShow = totalFeedbacks
	}

	// build feedback items HTML
	for i := 0; i < maxToShow; i++ {
		item := feedbacks[i]
		itemContent := r.feedbackItemTemplate

		// set project name
		itemContent = strings.Replace(itemContent, "{{projectName}}", item.Project.Title, 1)

		// handle user identification (app user id > email > Anonymous)
		var userIdentification string
		if *item.AppUserID != "" {
			userIdentification = fmt.Sprintf("<span style=\"color: #4b4e68;\">User: %s</span>", *item.AppUserID)
		} else if *item.Email != "" {
			userIdentification = fmt.Sprintf("<span style=\"color: #4b4e68;\">%s</span>", *item.Email)
		} else {
			userIdentification = "<span style=\"font-style: italic; color: #4b4e68;\">Anonymous</span>"
		}
		itemContent = strings.Replace(itemContent, "{{userIdentification}}", userIdentification, 1)

		// format created time with year
		createdAt := item.CreatedAt.Format("Jan 2, 2006 15:04")
		itemContent = strings.Replace(itemContent, "{{createdAt}}", createdAt, 1)

		// set feedback type and badge color
		itemContent = strings.Replace(itemContent, "{{feedbackType}}", strings.ToUpper(item.CampaignType.String()), 1)
		itemContent = strings.Replace(itemContent, "{{typeBadgeColor}}", "#c0c8d8", 1)
		itemContent = strings.Replace(itemContent, "{{typeBadgeTextColor}}", "#272f3f", 1)

		// set rating and max rating
		var maxRating int32
		if item.CampaignType.IsNPS() || item.CampaignType.IsCSAT() {
			maxRating = 10
		} else {
			maxRating = 5
		}

		// ensure rating is within bounds
		rating := item.Rating
		if rating <= 0 {
			rating = 0
		}
		if rating > maxRating {
			rating = maxRating
		}

		// set rating values
		itemContent = strings.Replace(itemContent, "{{rating}}", fmt.Sprintf("%d", rating), 1)
		itemContent = strings.Replace(itemContent, "{{maxRating}}", fmt.Sprintf("%d", maxRating), 1)

		// handle content section
		if item.Content != "" {
			itemContent = strings.Replace(itemContent, "{{#if hasContent}}", "", 1)
			itemContent = strings.Replace(itemContent, "{{/if}}", "", 1)
			itemContent = strings.Replace(itemContent, "{{content}}", item.Content, 1)
		} else {
			// Remove content section if no content
			itemContent = strings.Replace(itemContent, "{{#if hasContent}}", "<!-- ", 1)
			itemContent = strings.Replace(itemContent, "{{/if}}", " -->", 1)
		}

		feedbackItemsContent += itemContent
	}

	// replace placeholders in the main template
	content = strings.Replace(content, "{{feedbackCount}}", fmt.Sprintf("%d", totalFeedbacks), 1)
	content = strings.Replace(content, "{{feedbackLabel}}", feedbackLabel, 1)
	content = strings.Replace(content, "{{feedbackItemsContent}}", feedbackItemsContent, 1)
	content = strings.Replace(content, "{{dashboardUrl}}", "https://easyfeedback.com", 1)

	return r.mailer.SendEmail(ctx, mailer.SendEmailRequest{
		To:      toEmail,
		Subject: subject,
		Content: content,
	})
}
