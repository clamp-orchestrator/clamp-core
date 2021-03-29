package services

import (
	"clamp-core/repository"
	"os"
	"testing"
)

var mockDB repository.MockDB

func TestMain(m *testing.M) {
	repository.SetDB(&mockDB)
	os.Exit(m.Run())
}
