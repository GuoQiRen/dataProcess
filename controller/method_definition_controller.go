package controller

import (
	"dataProcess/service/local/impl"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"github.com/gin-gonic/gin"
)

func MethodDefinitionUpload(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	upMethodOrm, err := impl.UPMethodOrmDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	id, err := impl.CreateDefinitionImpl().UploadDefinitionMethod(upMethodOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]int32)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func MethodDefinitionUpdate(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id, err := utils.StringToint32(c.DefaultQuery("id", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	userId, err := utils.StringToint32(c.DefaultQuery("userId", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	upMethodOrm, err := impl.UPMethodOrmDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateDefinitionImpl().UploadDefinitionMethodUpdate(id, userId, upMethodOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]int32)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func MethodDefinitionDownload(c *gin.Context) {

}

func MethodDefinitionAuth(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")

	authInfos, err := impl.AuthInfoDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateDefinitionImpl().AssignDefinitionMethodAuthority(id, userId, authInfos.AuthInfo)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]string)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func MethodDefinitionAuthSelect(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")

	authInfos, err := impl.CreateDefinitionImpl().SelectDefinitionMethodAuthority(id, userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]interface{})
	ret["authInfo"] = authInfos

	app.OK(c, ret, true, "", 0)
}

func MethodDefinitionDelete(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")

	err = impl.CreateDefinitionImpl().DeleteDefinitionMethod(id, userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]string)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}
