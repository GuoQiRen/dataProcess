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

type OldTokenContext struct {
	Token string `json:"token"`
}

func CreateOldTokenContext(token string) *OldTokenContext {
	return &OldTokenContext{Token: token}
}

func (o *OldTokenContext) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	oldTokenContextBytes, err := json.Marshal(*o)
	if err != nil {
		return
	}
	oldTokenContextBody := strings.NewReader(string(oldTokenContextBytes))
	// 刷新token
	resp, err := app.UriRequest(request.POST, dbDo.DataSetConfig.Head+dbDo.DataSetConfig.Host+constants.Colon+dbDo.DataSetConfig.Port+
		dbDo.DataSetConfig.UserTokenContext, oldTokenContextBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	userTokenResp := mocks.CreateUserTokenResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&userTokenResp)
	if err != nil {
		return
	}

	return userTokenResp, err
}
