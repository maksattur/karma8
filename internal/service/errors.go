package service

import "errors"

var (
	ErrForeignKeyOrUniqueViolation = errors.New("foreign key or unique violation")
)
