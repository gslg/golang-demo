package util

import (
	"log"
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()

	if uuid == "" {
		t.Error("生成uuid报错")
	}

	log.Println(uuid)
}
