package impl

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"net/url"
)

type UserGroupsReq struct {
	PlatformCode string `json:"platformCode"`
	GroupId      int32  `json:"groupId"`
}

func (u *UserGroupsReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	query.Set("platformCode", u.PlatformCode)
	query.Set("groupId", utils.IntToString(int(u.GroupId)))
	resp, err := app.UriRequest(request.GET, config.TestPlatConfig.Head+config.TestPlatConfig.Host+constants.Colon+config.TestPlatConfig.Port+config.TestPlatConfig.UserGroupContext, nil, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	userGroupResp := mocks.CreateUserGroupResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&userGroupResp)
	if err != nil {
		return
	}

	return userGroupResp, err
}
