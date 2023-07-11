package errorhandler

type Error int

const (
	InternealServer Error = iota
	NotFound
	Forbidden
	AlreadyExist
)
