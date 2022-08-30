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

type ContainerConfig struct {
	GPU uint32 `json:"GPU"`
	CPU uint32 `json:"CPU"`
}

type ResourceGroup struct {
	ResGroupName           string  `json:"resGroupName"`
	ResGroupLabel          string  `json:"resGroupLabel"`
	ResourceGroupCloudId   string  `json:"resource_group_cloud_id"`
	ResourceGroupServerCpu int32   `json:"resource_group_server_cpu"`
	ResourceCpuOverRation  float64 `json:"resource_cpu_over_ration"`
	ResourceGroupType      string  `json:"resource_group_type"`
}

type TrainContainerReqContext struct {
	UserName             string           `json:"userName"`
	ImageName            string           `json:"imageName"`
	ImageTag             string           `json:"imageTag"`
	WorkPath             string           `json:"workPath"`
	ContainerConfig      *ContainerConfig `json:"containerConfig"`
	Dataset              []interface{}    `json:"dataset"`
	UserGroup            string           `json:"userGroup"`
	ResourceGroup        *[]ResourceGroup `json:"resourceGroup"`
	IsMulti              string           `json:"isMulti"`
	IsLoop               bool             `json:"isLoop"`
	NodeToUse            int32            `json:"nodeToUse"`
	Priority             int32            `json:"priority"`
	SecondStorePathList  []string         `json:"secondStorePathList"`
	GpuType              int32            `json:"gpuType"`
	HardwareType         int32            `json:"hardwareType"`
	Tag                  string           `json:"tag"`
	MissionName          string           `json:"missionName"`
	FrameType            string           `json:"frameType"`
	DatasetVersionIdList []interface{}    `json:"datasetVersionIdList"`
	ResourceModel        int32            `json:"resourceModel"`
}

func (d *TrainContainerReqContext) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	trainBytes, err := json.Marshal(*d)
	trainBody := strings.NewReader(string(trainBytes))

	// 启动dug容器任务
	resp, err := app.UriRequest(request.POST, config.TrainPlatConfig.Head+config.TrainPlatConfig.Host+config.TrainPlatConfig.JobContext, trainBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	jobResp := mocks.CreateJobResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jobResp)
	if err != nil {
		return
	}

	return jobResp, err
}
