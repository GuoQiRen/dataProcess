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

type FormalCreateDataset struct {
	VersionId int32 `json:"version_id"`
	ParentId  int32 `json:"parent_id"`
	Creator   int32 `json:"creator"`
}

func CreateFormalCreateDataset(creator int32, versionId int32) *FormalCreateDataset {
	return &FormalCreateDataset{VersionId: versionId, Creator: creator}
}

func (f *FormalCreateDataset) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	formalCreateBytes, err := json.Marshal(*f)
	formalDatasetBody := strings.NewReader(string(formalCreateBytes))

	// 发送预创建Dataset请求
	resp, err := app.UriRequest(request.PUT, dbDo.DataSetConfig.Head+dbDo.DataSetConfig.Host+constants.Colon+dbDo.DataSetConfig.Port+
		dbDo.DataSetConfig.FormalCreateDatasetContext, formalDatasetBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	formalDatasetResp := mocks.CreateFormalDatasetResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&formalDatasetResp)
	if err != nil {
		return
	}

	return formalDatasetResp, err
}
