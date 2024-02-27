package model

import "github.com/google/uuid"

var EmptyUUID = uuid.UUID{}

type Task struct {
	Description string
	Status      string
	Assignee    string
	ID          uuid.UUID
}

func (t Task) IsValidID() bool {
	if t.ID == EmptyUUID {
		return false
	}

	return true
}
