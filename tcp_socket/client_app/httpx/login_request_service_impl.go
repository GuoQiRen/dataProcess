package httpx

import (
	"dataProcess/constants"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	dbDo "dataProcess/tcp_socket/db_do"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
	"strings"
)

type LoginReq struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
	UserDomain  string `json:"userDomain"`
}

func CreateLoginReq(userAccount, password, userDomain string) LoginReq {
	return LoginReq{userAccount, password, userDomain}
}

func (l *LoginReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	loginBytes, err := json.Marshal(*l)
	if err != nil {
		return
	}
	loginBody := strings.NewReader(string(loginBytes))

	// 登录认证
	resp, err := app.UriRequest(request.POST, dbDo.DataSetConfig.Head+dbDo.DataSetConfig.Host+constants.Colon+dbDo.DataSetConfig.Port+dbDo.DataSetConfig.LoginContext, loginBody, query, heads)
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
