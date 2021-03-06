package types

import "fmt"

// Error represents a generic Hephaestus error
type Error struct {
	Description string
	Reaction    string
}

// New allows to build a new Error instance
func New(description, reaction string) *Error {
	return &Error{
		Description: description,
		Reaction:    reaction,
	}
}

// Error implements err.Error
func (e *Error) Error() string {
	return e.Description
}

// -------------------------------------------------------------------------------------------------------------------

func NewWarnErr(description string, args ...interface{}) *Error {
	return New(fmt.Sprintf(description, args...), "⚠️")
}

func NewTimeErr(description string) *Error {
	return New(description, "⌛")
}
