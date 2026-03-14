package valueobject

import "errors"

type BookStatus string

const (
	Unread   BookStatus = "unread"
	Reading  BookStatus = "reading"
	Finished BookStatus = "finished"
)

var ErrInvalidBookStatus = errors.New("invalid book status")

func NewBookStatus(s string) (BookStatus, error) {
	switch BookStatus(s) {
	case Unread, Reading, Finished:
		return BookStatus(s), nil
	default:
		return "", ErrInvalidBookStatus
	}
}

func (bs BookStatus) String() string {
	return string(bs)
}

func (bs BookStatus) CanTransitionTo(next BookStatus) bool {
	switch bs {
	case Unread:
		return next == Reading
	case Reading:
		return next == Finished
	case Finished:
		return false
	default:
		return false
	}
}
