package state

import "fmt"

type StateReadAlert struct {
	key string
	err string
}

func NewStateReadAlert(key string, err error) *StateReadAlert {
	return &StateReadAlert{key: key, err: err.Error()}
}

func (s *StateReadAlert) Message() string {
	return fmt.Sprintf("Your analysis will be incomplete. There was an error reading state file '%s': %s", s.key, s.err)
}

func (s *StateReadAlert) ShouldIgnoreResource() bool {
	return false
}

type StateEnumerationAlert struct {
	source string
	err    string
}

func NewStateEnumerationAlertAlert(err error) *StateEnumerationAlert {
	return &StateEnumerationAlert{err: err.Error()}
}

func (s *StateEnumerationAlert) Message() string {
	return fmt.Sprintf("Your analysis will be incomplete. There was an error enumerating statefile for '%s': %s", s.source, s.err)
}

func (s *StateEnumerationAlert) ShouldIgnoreResource() bool {
	return false
}

type StateEnumerationError struct {
}
