package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user already exists")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTaskStatus = errors.New("invalid task status")
) 