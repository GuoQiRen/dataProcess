package impl

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
)

type UserDetailReq struct {
	UserAccount  string `json:"userAccount"`
	PlatFormCode string `json:"platformCode"`
}

func (u *UserDetailReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	query.Set("userAccount", u.UserAccount)
	query.Set("platformCode", u.PlatFormCode)
	resp, err := app.UriRequest(request.GET, config.TestPlatConfig.Head+config.TestPlatConfig.Host+constants.Colon+config.TestPlatConfig.Port+config.TestPlatConfig.UserDetailContext, nil, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	userDetailResp := mocks.CreateUserDetailResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&userDetailResp)
	if err != nil {
		return
	}

	return userDetailResp, err
}
