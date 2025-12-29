package storage

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrTaskExists = errors.New("task exists")
	ErrEmptyTitle = errors.New("empty task title")
	ErrTaskAlreadyComplited = errors.New("task already complited")
)