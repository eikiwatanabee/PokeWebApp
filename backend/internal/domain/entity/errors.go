package entity

import "errors"

var (
	ErrEmptyTitle          = errors.New("book title cannot be empty")
	ErrEmptyAuthor         = errors.New("book author cannot be empty")
	ErrInvalidTransition   = errors.New("invalid status transition")
	ErrAlreadyFinished     = errors.New("book is already finished")
	ErrNotReading          = errors.New("book is not in reading status")
	ErrEmptyMemoContent    = errors.New("memo content cannot be empty")
	ErrEmptyTagName        = errors.New("tag name cannot be empty")
	ErrEmptyTenantName     = errors.New("tenant name cannot be empty")
	ErrEmptyEmail          = errors.New("email cannot be empty")
	ErrEmptyGoogleID       = errors.New("google ID cannot be empty")
	ErrEmptyGitHubID       = errors.New("github ID cannot be empty")
	ErrEmptyTeamName       = errors.New("team name cannot be empty")
)
