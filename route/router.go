package route

import (
	"github.com/gin-gonic/gin"
	"secKill/controller"
)

func InitRouter() (router *gin.Engine) {
	router = gin.Default()
	router.POST("/secKill", controller.SecKill)
	router.POST("/secKillInfo", controller.GetProductInfo)
	return
}
