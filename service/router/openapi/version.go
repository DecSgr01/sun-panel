package openapi

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/middleware"
	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
)

func InitVersion(router *gin.RouterGroup) {
	r := router.Group("", middleware.PublicTokenInterceptor)
	// 直接返回1.7.0版本号
	r.POST("/v1/version", func(c *gin.Context) {
		version := cmn.GetSysVersionInfo()
		apiReturn.SuccessData(c, gin.H{
			"version":     "1.7.0",
			"versionCode": version.Version_code,
		})
	})
}
