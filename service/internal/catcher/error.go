package catcher

// Error is our custom error
type Error string

const (
	ErrProducerStopped = Error("producer stopped")
	ErrNotFound = Error("not found")
	ErrInvalid = Error("invalid")
)

// Error implements the error interface
func (e Error) Error() string {
	return string(e)
}


