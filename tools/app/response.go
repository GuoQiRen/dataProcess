package app

import (
	"dataProcess/logger"
	"dataProcess/repository/mocks"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

/**
初始化response实体
*/
func initVar(result bool, errMsg string, code int32) *ReturnData {
	var ret ReturnData
	returnData := ret.CreateReturn(result, errMsg, code)
	return returnData
}

/**
返回error
*/
func Error(c *gin.Context, result bool, errMsg string, code int32) {
	logger.Error(errMsg)
	c.JSON(http.StatusOK, initVar(result, errMsg, code))
}

/**
返回success
*/
func OK(c *gin.Context, data interface{}, result bool, errMsg string, code int32) {
	returnData := initVar(result, errMsg, code)
	returnData.Data = data
	c.JSON(http.StatusOK, returnData)
}

/**
分页
*/
func SlicePage(pageNum, pageSize, nums int32) (sliceStart, sliceEnd int32) {
	if pageNum < 0 {
		pageNum = 1
	}

	if pageSize < 0 {
		pageSize = 12
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int32(math.Ceil(float64(nums) / float64(pageSize)))
	if pageNum > pageCount {
		return 0, 0
	}
	sliceStart = (pageNum - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

/**
分页成功
*/
func PageOK(c *gin.Context, data interface{}, result bool, errMsg string, total int32, code int32) {
	var pageR PageResponse
	pageResp := pageR.CreatePageResponse(result, errMsg, code, total)
	pageResp.Data = data
	c.JSON(http.StatusOK, pageResp)
}

/**
解析page参数
*/
func PageDecode(c *gin.Context) (pageCore mocks.PageCore, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&pageCore)
	if err != nil {
		return
	}
	return
}
