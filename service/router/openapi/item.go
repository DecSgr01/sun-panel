package openapi

import (
	"sun-panel/api/api_v1/middleware"
	"sun-panel/api/openapi"

	"github.com/gin-gonic/gin"
)

func InitIiem(router *gin.RouterGroup) {
	item := openapi.ApiApp.Apiv1.Item
	r := router.Group("", middleware.PublicTokenInterceptor)
	r.POST("/v1/item/create", item.Create)
}
