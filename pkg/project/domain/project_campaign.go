package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectCampaignRepository interface {
	Create(ctx *appcontext.AppContext, campaign ProjectCampaign) error
	Update(ctx *appcontext.AppContext, campaign ProjectCampaign) error
	FindByID(ctx *appcontext.AppContext, campaignID string) (*ProjectCampaign, error)
	FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]ProjectCampaign, error)
}

type ProjectCampaign struct {
	ID                    string
	ProjectID             string
	Name                  string
	Description           string
	CampaignType          ProjectCampaignType
	Status                Status
	SettingWidgetPosition string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func NewProjectCampaign() (*ProjectCampaign, error) {
	var (
		now = manipulation.NowUTC()
	)

	var c = &ProjectCampaign{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return c, nil
}

func (c *ProjectCampaign) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}
	c.ProjectID = projectID
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCampaign) SetName(name string) error {
	if name == "" || len(name) < 3 || len(name) > 255 {
		return apperrors.Common.InvalidName
	}

	c.Name = name
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCampaign) SetDescription(description string) error {
	if len(description) > 2000 {
		return apperrors.Common.InvalidDescription
	}

	c.Description = description
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCampaign) SetCampaignType(campaignType string) error {
	var dCampaignType = ToProjectCampaignType(campaignType)
	if !dCampaignType.IsValid() {
		return apperrors.Project.InvalidCampaignType
	}

	c.CampaignType = dCampaignType
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCampaign) SetStatus(status string) error {
	var dStatus = ToStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	c.Status = dStatus
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCampaign) SetSettingWidgetPosition(settingWidgetPosition string) error {
	if settingWidgetPosition != "" {
		c.SettingWidgetPosition = settingWidgetPosition
		c.SetUpdatedAt()
	}

	return nil
}

func (c *ProjectCampaign) SetUpdatedAt() {
	c.UpdatedAt = manipulation.NowUTC()
}
