package util

import (
	"github.com/satori/go.uuid"
	"strings"
)

func NewUUID() string{
	return strings.Replace(uuid.Must(uuid.NewV4()).String(),"-","",-1)
}