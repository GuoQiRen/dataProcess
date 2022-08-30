package impl

import (
	"dataProcess/config"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
)

type TrainLogReq struct {
	FilePath string `json:"filePath"`
}

func (e *TrainLogReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	query.Set("filePath", e.FilePath)
	resp, err := app.UriRequest(request.GET, config.TrainPlatConfig.Head+config.TrainPlatConfig.Host+config.TrainPlatConfig.TrainLogContext, nil, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	taskTrainLogResp := mocks.CreateTaskTrainLogResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&taskTrainLogResp)
	if err != nil {
		return
	}

	return taskTrainLogResp, err
}
