package shared

import (
	"uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.NewV7())
}

func (id ID) String() string {
	return id.String()
}

func ParseID(raw string) (ID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ID(uuid.Nil()), ErrIDParsing
	}
	return ID(parsed), nil
}
