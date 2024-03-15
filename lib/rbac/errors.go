package rbac

import (
	"errors"
)

var (
	ErrUidInvalid  = errors.New("uid is invalid")
	ErrGidInvalid  = errors.New("gid is invalid")
	ErrPermInvalid = errors.New("perm is invalid")

	ErrUidAlreadyExists = errors.New("uid already exists")
	ErrGidAlreadyExists = errors.New("gid already exists")

	ErrUidNotFound = errors.New("uid not found")
	ErrGidNotFound = errors.New("gid not found")
)
