package impl

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/request"
	"dataProcess/repository/mocks"
	"dataProcess/tools/app"
	"encoding/json"
	"net/url"
	"strconv"
)

type DataSetReqContext struct {
	SolidId int32 `json:"solidId"`
}

func (d *DataSetReqContext) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	// 查询数据集的真实id
	query.Set("ver_id", strconv.Itoa(int(d.SolidId)))
	resp, err := app.UriRequest(request.GET, config.TestPlatConfig.Head+config.TestPlatConfig.Host+constants.Colon+config.TestPlatConfig.Port+config.TestPlatConfig.DataSetContext, nil, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dataSetResp := mocks.CreateDataSetResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dataSetResp)
	if err != nil {
		return
	}

	return dataSetResp, err
}
