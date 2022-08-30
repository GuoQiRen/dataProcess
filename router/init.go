package router

import (
	"dataProcess/router/system"
	"github.com/gin-gonic/gin"
)

func InitSysRouter(r *gin.Engine) *gin.RouterGroup {
	g := r.Group("")

	// 无需认证
	system.SysRouterInit(g)

	return g
}

func InitRouter() *gin.Engine {
	r := gin.New()

	// 注册系统路由
	InitSysRouter(r)

	return r
}
