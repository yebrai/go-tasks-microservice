package task

import (
	"errors"

	"github.com/yebrai/go-tasks-microservice/pkg/id"
)

type ID string

var ErrInvalidID = errors.New("invalid ID")

func NewID(value string) (ID, error) {
	if !id.IsValidUID(value) {
		return "", ErrInvalidID
	}
	return ID(value), nil
}
