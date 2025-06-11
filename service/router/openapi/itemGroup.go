package openapi

import (
	"sun-panel/api/api_v1/middleware"
	"sun-panel/api/openapi"

	"github.com/gin-gonic/gin"
)

func InitIiemGroup(router *gin.RouterGroup) {
	itemGroup := openapi.ApiApp.Apiv1.ItemGroup
	r := router.Group("", middleware.PublicTokenInterceptor)
	r.POST("/v1/itemGroup/getList",itemGroup.GetList)
}
