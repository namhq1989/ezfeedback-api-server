package apperrors

import "errors"

const prefixUniqueViolation = "unique_violation_"

func GetErrorUniqueViolationByConstraintName(constraintName string) error {
	return errors.New(prefixUniqueViolation + constraintName)
}
