package initialize

import (
	"cdbot/src/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() (*gin.Engine, error) {
	router := gin.Default()

	router.POST("/", api.GenericMessageHandler)

	return router, nil
}
