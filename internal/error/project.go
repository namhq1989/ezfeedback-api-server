package apperrors

import "errors"

var Project = struct {
	InvalidProjectID      error
	ProjectNotFound       error
	InvalidRole           error
	InvalidCategory       error
	CategoryLimitExceeded error
	InvalidCampaign       error
	InvalidCampaignType   error
	InvalidDomain         error
	InvalidPrimaryColor   error
}{
	InvalidProjectID:      errors.New("project_invalid_id"),
	ProjectNotFound:       errors.New("project_not_found"),
	InvalidRole:           errors.New("project_invalid_role"),
	InvalidCategory:       errors.New("project_invalid_category"),
	CategoryLimitExceeded: errors.New("project_category_limit_exceeded"),
	InvalidCampaign:       errors.New("project_invalid_campaign"),
	InvalidCampaignType:   errors.New("project_invalid_campaign_type"),
	InvalidDomain:         errors.New("project_invalid_domain"),
	InvalidPrimaryColor:   errors.New("project_invalid_primary_color"),
}
