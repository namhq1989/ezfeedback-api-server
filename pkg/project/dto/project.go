package dto

import (
	"slices"

	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectBrief struct {
	ID     string       `json:"id"`
	Title  string       `json:"title"`
	Slug   string       `json:"slug"`
	Status string       `json:"status"`
	Stats  ProjectStats `json:"stats"`
}

func (ProjectBrief) FromDomain(project domain.Project) ProjectBrief {
	return ProjectBrief{
		ID:     project.ID,
		Title:  project.Title,
		Slug:   project.Slug,
		Status: project.Status.String(),
		Stats: ProjectStats{
			TotalFeedbacks: project.StatsTotalFeedbacks,
		},
	}
}

type Project struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Slug        string                    `json:"slug"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	Campaigns   []ProjectCampaign         `json:"campaigns"`
	Categories  []ProjectCategory         `json:"categories"`
	Setting     ProjectSetting            `json:"setting"`
	Stats       ProjectStats              `json:"stats"`
	CreatedAt   *httprespond.TimeResponse `json:"createdAt"`
	UpdatedAt   *httprespond.TimeResponse `json:"updatedAt"`
}

type ProjectStats struct {
	TotalFeedbacks int32 `json:"totalFeedbacks"`
}

func (Project) FromDomain(
	project domain.Project,
	setting domain.ProjectSetting,
	campaigns []domain.ProjectCampaign,
	categories []domain.ProjectCategory,
) Project {
	var cats = make([]ProjectCategory, len(categories))
	for i, category := range categories {
		cats[i] = ProjectCategory{}.FromDomain(category)
	}

	var campns = make([]ProjectCampaign, len(campaigns))
	for i, campaign := range campaigns {
		campaignCategories := make([]domain.ProjectCategory, 0)
		for _, category := range categories {
			if slices.Contains(campaign.Categories, category.ID) {
				campaignCategories = append(campaignCategories, category)
			}
		}

		campns[i] = ProjectCampaign{}.FromDomain(campaign, campaignCategories)
	}

	return Project{
		ID:          project.ID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description,
		Status:      project.Status.String(),
		Categories:  cats,
		Campaigns:   campns,
		Setting:     ProjectSetting{}.FromDomain(setting),
		Stats: ProjectStats{
			TotalFeedbacks: project.StatsTotalFeedbacks,
		},
		CreatedAt: httprespond.NewTimeResponse(project.CreatedAt),
		UpdatedAt: httprespond.NewTimeResponse(project.UpdatedAt),
	}
}
