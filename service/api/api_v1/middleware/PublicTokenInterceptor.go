package middleware

import (
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func PublicTokenInterceptor(c *gin.Context) {
	// 获得token
	username := c.GetHeader("token")

	// 没有token信息视为未登录
	if username != "" {
		global.Logger.Debug("数据库查询username:", username)
		mUser := models.User{}
		// 去库中查询是否存在该用户
		if info, err := mUser.GetUserInfoByUsername(username); err == nil && info.Token != "" && info.ID != 0 {
			global.Logger.Debug("数据库查询用户:", info.ID)
			c.Set("userInfo", info)
			return
		} else {
			global.Logger.Debug("数据库查询用户失败", username)
		}
	} else {
		global.Logger.Debug("username不存在")
	}
	c.Abort()
}
