package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type Project struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Slug        string                    `json:"slug"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	Categories  []ProjectCategory         `json:"categories"`
	Setting     ProjectSetting            `json:"setting"`
	Stats       ProjectStats              `json:"stats"`
	CreatedAt   *httprespond.TimeResponse `json:"createdAt"`
	UpdatedAt   *httprespond.TimeResponse `json:"updatedAt"`
}

type ProjectStats struct {
	TotalFeedback int `json:"totalFeedback"`
}

func (Project) FromDomain(project domain.Project, setting domain.ProjectSetting, categories []domain.ProjectCategory) Project {
	var cats = make([]ProjectCategory, len(categories))
	for i, category := range categories {
		cats[i] = ProjectCategory{}.FromDomain(category)
	}

	return Project{
		ID:          project.ID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description,
		Status:      project.Status.String(),
		Categories:  cats,
		Setting:     ProjectSetting{}.FromDomain(setting),
		Stats: ProjectStats{
			TotalFeedback: 10,
		},
		CreatedAt: httprespond.NewTimeResponse(project.CreatedAt),
		UpdatedAt: httprespond.NewTimeResponse(project.UpdatedAt),
	}
}
