package shared

import (
	"uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.NewV7())
}
