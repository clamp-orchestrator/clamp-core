package handlers

import (
	"clamp-core/config"
	_ "clamp-core/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"time"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.ENV.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/metrics", prometheusHandler())
	r.POST("/serviceRequest/:workflowName", createServiceRequestHandler())
	r.GET("/serviceRequest/:serviceRequestId", getServiceRequestStatusHandler())
	r.POST("/workflow", createWorkflowHandler())
	r.GET("/workflow/:workflow", fetchWorkflowBasedOnWorkflowName())
	r.POST("/stepResponse", createStepResponseHandler())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/workflows", getWorkflows())
	r.GET("/serviceRequests/:workflowName", findServiceRequestByWorkflowNameHandler())
	return r
}

func LoadHTTPRoutes() {
	r := setupRouter()
	err := r.Run()
	panic(err)
}
