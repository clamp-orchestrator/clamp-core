package handlers

import (
	"clamp-core/repository"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var mockDB repository.MockDB

func TestMain(m *testing.M) {
	repository.SetDB(repository.NewMemoryDB())
	gin.SetMode(gin.TestMode)
	//log.SetLevel(log.DebugLevel)
	os.Exit(m.Run())
}
