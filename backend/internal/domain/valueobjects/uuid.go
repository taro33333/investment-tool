package valueobjects

import (
	"github.com/google/uuid"
)

type UUID struct {
	value string
}

func NewUUID() UUID {
	return UUID{value: uuid.New().String()}
}

func ParseUUID(s string) (UUID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return UUID{}, err
	}
	return UUID{value: s}, nil
}

func (id UUID) String() string {
	return id.value
}

func (id UUID) Equals(other UUID) bool {
	return id.value == other.value
}

func (id UUID) IsZero() bool {
	return id.value == ""
}
