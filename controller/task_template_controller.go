package controller

import (
	"dataProcess/service/local/impl"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"github.com/gin-gonic/gin"
)

func TaskTemplateSave(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	templateOrm, err := impl.TemplateCompletePostDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	id, err := impl.CreateTaskTemplateImpl().SaveTemplate(templateOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]int32)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func TaskTemplateUpdate(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	templateOrm, err := impl.TemplateCompletePostDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateTaskTemplateImpl().UpdateTemplate(templateOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]int32)
	ret["id"] = templateOrm.Id

	app.OK(c, ret, true, "", 0)
}

func TaskTemplateSelect(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	idI, err := utils.StringToint32(c.DefaultQuery("id", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	template, err := impl.CreateTaskTemplateImpl().SelectTemplate(idI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	app.OK(c, template, true, "", 0)
}

func TaskTemplateDelete(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	idI, err := utils.StringToint32(c.DefaultQuery("id", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	userId, err := utils.StringToint32(c.DefaultQuery("userId", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateTaskTemplateImpl().DeleteTemplate(idI, userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, "", true, "", 0)
}

func TaskTemplateSelectAll(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	userId, err := utils.StringToint32(c.DefaultQuery("userId", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	templates, err := impl.CreateTaskTemplateImpl().SelectAllTemplate(userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, templates, true, "", 0)
}
