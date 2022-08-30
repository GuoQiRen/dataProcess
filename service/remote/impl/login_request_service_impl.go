package impl

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
	"strings"
)

type LoginReqContext struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
	UserDomain  string `json:"userDomain"`
}

func (l *LoginReqContext) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	loginBytes, err := json.Marshal(*l)
	if err != nil {
		return
	}
	loginBody := strings.NewReader(string(loginBytes))

	// 登录认证
	resp, err := app.UriRequest(request.POST, config.TestPlatConfig.Head+config.TestPlatConfig.Host+constants.Colon+config.TestPlatConfig.Port+config.TestPlatConfig.LoginContext, loginBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	loginResp := mocks.CreateLoginResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&loginResp)
	if err != nil {
		return
	}
	return loginResp, err
}
