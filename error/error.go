package error

import (
	"errors"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrForbidden            = errors.New("forbidden")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrProjectNotFound      = errors.New("project not found")
	ErrNotProjectOwner      = errors.New("only the project owner can perform this action")
	ErrAlreadyProjectMember = errors.New("user is already a member of this project")
	ErrNotProjectMember     = errors.New("user is not a member of this project")
	ErrInvalidInput         = errors.New("invalid input")
	ErrConflict             = errors.New("resource conflict")
	ErrInternal             = errors.New("internal server error")
	ErrBadRequest           = errors.New("bad request")
	ErrTokenExpired         = errors.New("token expired")
)
