package controller

import (
	"dataProcess/service/local/impl"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"github.com/gin-gonic/gin"
)

func TaskPrototypeSave(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	protoOrm, err := impl.PrototypeDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	protoType, err := impl.CreateTaskPrototypeImpl().SaveProto(protoOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	app.OK(c, protoType, true, "", 0)
}

func TaskPrototypeUpdate(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	protoOrm, err := impl.PrototypeDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	err = impl.CreateTaskPrototypeImpl().UpdateProto(protoOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, protoOrm, true, "", 0)
}

func TaskPrototypeSelect(c *gin.Context) {
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

	protoType, err := impl.CreateTaskPrototypeImpl().SelectProto(idI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	app.OK(c, protoType, true, "", 0)
}

func TaskPrototypeDelete(c *gin.Context) {
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

	err = impl.CreateTaskPrototypeImpl().DeleteProto(idI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	app.OK(c, "", true, "", 0)
}
