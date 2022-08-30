package system

import (
	"dataProcess/controller"
	"github.com/gin-gonic/gin"
)

func SysRouterInit(r *gin.RouterGroup) {
	process := r.Group("/dldmp/dataprocess")
	{
		// 方法类型管理
		process.POST("/method/category/create", controller.MethodCategoryCreate)
		process.PUT("/method/category/update", controller.MethodCategoryUpdate)
		process.POST("/method/category/select/ids", controller.MethodCategorySelectByIds)
		process.GET("/method/category/select/parentId", controller.MethodCategorySelectByParentId)
		process.GET("/method/category/select/all", controller.MethodCategoryAll)
		process.DELETE("/method/category/delete", controller.MethodCategoryDelete)

		// 方法上传管理
		process.POST("/method/define/upload", controller.MethodDefinitionUpload)
		process.PUT("/method/define/update", controller.MethodDefinitionUpdate)
		process.PUT("/method/define/auth", controller.MethodDefinitionAuth)
		process.GET("/method/define/auth/select", controller.MethodDefinitionAuthSelect)
		process.DELETE("/method/define/delete", controller.MethodDefinitionDelete)

		// 任务管理
		process.GET("/task/log", controller.TaskManageLog)
		process.GET("/task/copy", controller.TaskManageCopy)
		process.POST("/task/save", controller.TaskManageSave)
		process.PUT("/task/update", controller.TaskManageUpdate)
		process.GET("/task/start", controller.TaskManageStart)
		process.POST("/task/search", controller.TaskManageSearch)
		process.GET("/task/detail", controller.TaskManageDetail)
		process.DELETE("/task/delete", controller.TaskManageDelete)
		process.GET("/task/terminal", controller.TaskManageTerminal)

		// 任务原型管理
		process.POST("/task/proto/create", controller.TaskPrototypeSave)
		process.POST("/task/proto/update", controller.TaskPrototypeUpdate)
		process.GET("/task/proto/select", controller.TaskPrototypeSelect)
		process.DELETE("/task/proto/delete", controller.TaskPrototypeDelete)

		// 模版管理
		process.POST("/template/create", controller.TaskTemplateSave)
		process.PUT("/template/update", controller.TaskTemplateUpdate)
		process.GET("/template/select", controller.TaskTemplateSelect)
		process.DELETE("/template/delete", controller.TaskTemplateDelete)
		process.GET("/template/select/all", controller.TaskTemplateSelectAll)
	}
}
