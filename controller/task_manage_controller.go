package controller

import (
	"dataProcess/cache/redis"
	"dataProcess/constants"
	"dataProcess/service/local/impl"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"github.com/gin-gonic/gin"
)

func TaskManageSave(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	taskOrm, err := impl.TaskPostDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	id, err := impl.CreateTaskManageImpl().SaveTask(taskOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	ret := make(map[string]int32)
	ret["id"] = id
	app.OK(c, ret, true, "", 0)
}

func TaskManageStart(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")
	token, err := c.Cookie("token")
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = redis.SetOperate(constants.JinnToken+constants.SubLine+userId, token, -1)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateTaskManageImpl().StartTask(id, userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]string)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func TaskManageSearch(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	queryTask, err := impl.TaskGetDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	taskOrms, err := impl.CreateTaskManageImpl().SearchTask(&queryTask)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	start, end := app.SlicePage(queryTask.PageNum, queryTask.PageSize, int32(len(taskOrms)))

	app.PageOK(c, taskOrms[start:end], true, "", int32(len(taskOrms)), 0)
}

func TaskManageDetail(c *gin.Context) {
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

	taskOrm, err := impl.CreateTaskManageImpl().QueryTaskDetail(idI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, taskOrm, true, "", 0)
}

func TaskManageDelete(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")

	idI, err := utils.StringToint32(id)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	userIdI, err := utils.StringToint32(userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	taskM := impl.CreateTaskManageImpl()
	go func(taskM *impl.TaskManageImpl, id string) {
		_ = taskM.DeleteRemoteTask(id)
	}(taskM, id)

	err = taskM.DeleteTask(idI, userIdI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, "", true, "", 0)
}

func TaskManageUpdate(c *gin.Context) {
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

	taskOrm, err := impl.TaskPostDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateTaskManageImpl().UpdateTask(idI, taskOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	ret := make(map[string]int32)
	ret["id"] = idI
	app.OK(c, ret, true, "", 0)
}

func TaskManageTerminal(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	_ = c.DefaultQuery("userId", "")

	err = impl.CreateTaskManageImpl().TerminalTask(id)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]string)
	ret["id"] = id

	app.OK(c, ret, true, "", 0)
}

func TaskManageLog(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")

	logContent, err := impl.CreateTaskManageImpl().LogTask(id)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]string)
	ret["taskLog"] = logContent

	app.OK(c, ret, true, "", 0)
}

func TaskManageCopy(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	id := c.DefaultQuery("id", "")
	userId := c.DefaultQuery("userId", "")
	userName := c.DefaultQuery("userName", "")

	newId, err := impl.CreateTaskManageImpl().CopyTask(id, userId, userName)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	ret := make(map[string]int32)
	ret["id"] = newId

	app.OK(c, ret, true, "", 0)
}
