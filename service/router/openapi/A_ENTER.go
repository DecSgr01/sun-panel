package openapi

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitVersion(routerGroup)
	InitIiemGroup(routerGroup)
	InitIiem(routerGroup)
}
