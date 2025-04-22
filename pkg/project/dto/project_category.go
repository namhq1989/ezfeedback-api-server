package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

type ProjectCategory struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

func (ProjectCategory) FromDomain(projectCategory domain.ProjectCategory) ProjectCategory {
	return ProjectCategory{
		ID:     projectCategory.ID,
		Name:   projectCategory.Name,
		Slug:   projectCategory.Slug,
		Status: projectCategory.Status.String(),
	}
}
