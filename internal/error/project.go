package apperrors

import "errors"

var Project = struct {
	InvalidProjectID error
	ProjectNotFound  error
	InvalidRole      error
	InvalidCategory  error
}{
	InvalidProjectID: errors.New("project_invalid_id"),
	ProjectNotFound:  errors.New("project_not_found"),
	InvalidRole:      errors.New("project_invalid_role"),
	InvalidCategory:  errors.New("project_invalid_category"),
}
