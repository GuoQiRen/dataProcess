package impl

import (
	"dataProcess/config"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
)

type RuntimeReqContext struct {
}

func (e *RuntimeReqContext) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	// 查询runtime信息
	resp, err := app.UriRequest(request.GET, config.TrainPlatConfig.Head+config.TrainPlatConfig.Host+config.TrainPlatConfig.RuntimeContext, nil, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	runtimeResp := mocks.CreateRuntimeResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&runtimeResp)
	if err != nil {
		return
	}
	return runtimeResp, err
}
