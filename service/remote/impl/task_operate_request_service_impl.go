package impl

import (
	"dataProcess/config"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
	"strings"
)

type TaskOperateReq struct {
	UserName         string `json:"userName"`
	Operation        string `json:"operation"`
	MissionCreatedId string `json:"mission_created_id"`
	MissionStatus    string `json:"missionStatus"`
}

func (t *TaskOperateReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	taskOperateBytes, err := json.Marshal(*t)
	if err != nil {
		return
	}
	taskOperateBody := strings.NewReader(string(taskOperateBytes))

	// 进行任务操作
	resp, err := app.UriRequest(request.POST, config.TrainPlatConfig.Head+config.TrainPlatConfig.Host+config.TrainPlatConfig.TaskOperateContext, taskOperateBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	taskOperateResp := mocks.CreateTaskOperateResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&taskOperateResp)
	if err != nil {
		return
	}
	return taskOperateResp, err
}
