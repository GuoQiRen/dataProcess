package controller

import (
	"dataProcess/model"
	"dataProcess/service/local/impl"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"github.com/gin-gonic/gin"
)

func MethodCategoryCreate(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	cateOrm, err := impl.CategoryDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	err = impl.CreateCategoryImpl().CreateCategory(cateOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, cateOrm, true, "", 0)
}

func MethodCategoryUpdate(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	cateOrm, err := impl.CategoryDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateCategoryImpl().UpdateCategory(cateOrm)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, cateOrm, true, "", 0)
}

func MethodCategorySelectByIds(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	ids, err := impl.CategoryIdsDecode(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	categoryOrms, err := impl.CreateCategoryImpl().SelectCategoryByIds(ids)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, categoryOrms, true, "", 0)
}

func MethodCategorySelectByParentId(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	parentId, err := utils.StringToint32(c.DefaultQuery("parentId", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	categoryOrms, err := impl.CreateCategoryImpl().SelectCategoryByParentId(parentId)

	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, categoryOrms, true, "", 0)
}

func MethodCategoryAll(c *gin.Context) {
	err := TokenSecurityValid(c)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}
	name := c.DefaultQuery("name", "")
	userId := c.DefaultQuery("userId", "")

	init := model.CategoryTree{Children: make([]*model.CategoryTree, 0)}
	matchTree, err := impl.CreateCategoryImpl().SelectAllCategory(&init, name, userId)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	// 返回匹配节点数据
	if len(matchTree) > 0 {
		app.OK(c, matchTree, true, "", 0)
		return
	}

	app.OK(c, init, true, "", 0)
}

func MethodCategoryDelete(c *gin.Context) {
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

	userIdI, err := utils.StringToint32(c.DefaultQuery("userId", ""))
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	err = impl.CreateCategoryImpl().DeleteCategory(idI, userIdI)
	if err != nil {
		app.Error(c, false, err.Error(), -1)
		return
	}

	app.OK(c, "", true, "", 0)
}
