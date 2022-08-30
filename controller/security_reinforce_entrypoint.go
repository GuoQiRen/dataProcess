package controller

import (
	"dataProcess/security"
	"github.com/gin-gonic/gin"
)

var sessData security.SessData

func TokenSecurityValid(c *gin.Context) (err error) {
	token, err := c.Cookie("token")
	if err != nil {
		return
	}
	sessData.Token = token
	err = security.GetTokenInfo(&sessData)
	if err != nil {
		return
	}
	err = security.SetCookie(c, sessData.UserData.Token)
	if err != nil {
		return
	}
	return
}
