package domain

type ProjectCampaignType string

const (
	ProjectCampaignTypeUnknown  ProjectCampaignType = ""
	ProjectCampaignTypeFeedback ProjectCampaignType = "feedback"
	ProjectCampaignTypeNPS      ProjectCampaignType = "nps"
	ProjectCampaignTypeCSAT     ProjectCampaignType = "csat"
)

func (s ProjectCampaignType) String() string {
	return string(s)
}

func (s ProjectCampaignType) IsValid() bool {
	return s != ProjectCampaignTypeUnknown
}

func (s ProjectCampaignType) IsFeedback() bool {
	return s == ProjectCampaignTypeFeedback
}

func (s ProjectCampaignType) IsNPS() bool {
	return s == ProjectCampaignTypeNPS
}

func (s ProjectCampaignType) IsCSAT() bool {
	return s == ProjectCampaignTypeCSAT
}

func ToProjectCampaignType(s string) ProjectCampaignType {
	switch s {
	case ProjectCampaignTypeFeedback.String():
		return ProjectCampaignTypeFeedback
	case ProjectCampaignTypeNPS.String():
		return ProjectCampaignTypeNPS
	case ProjectCampaignTypeCSAT.String():
		return ProjectCampaignTypeCSAT
	default:
		return ProjectCampaignTypeUnknown
	}
}
