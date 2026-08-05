package apperrors

import "errors"

var (
	ErrInvalidCurrentPassword   = errors.New("invalid current password")
	ErrInvalidNewPassword       = errors.New("new password cannot be the same as current password")
	ErrInvalidPermissions       = errors.New("invalid permissions")
	ErrCannotDeleteSelf         = errors.New("cannot delete own user account")
	ErrCannotEditOwnPermissions = errors.New("cannot edit own permissions")
	ErrInvalidLoginUser         = errors.New("invalid login user")
	ErrInvalidLoginPassword     = errors.New("invalid login password")
)