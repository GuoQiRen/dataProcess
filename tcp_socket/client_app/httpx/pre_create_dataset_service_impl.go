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

type FileInfo struct {
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type PreCreateDataset struct {
	VersionName string     `json:"versionName"`
	Memo        string     `json:"tag"`
	SpaceId     string     `json:"spaceId"`
	FileInfos   []FileInfo `json:"fileInfos"`
	Creator     string     `json:"creator"`
}

func CreatePreCreateDataset(versionName, memo, spaceId, creator string, fileInfos []FileInfo) *PreCreateDataset {
	return &PreCreateDataset{VersionName: versionName, Memo: memo, SpaceId: spaceId, FileInfos: fileInfos, Creator: creator}
}

func (p *PreCreateDataset) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	preDatasetBytes, err := json.Marshal(*p)
	if err != nil {
		return
	}
	preDatasetBody := strings.NewReader(string(preDatasetBytes))

	// 发送预创建Dataset请求
	resp, err := app.UriRequest(request.POST, dbDo.DataSetConfig.Head+dbDo.DataSetConfig.Host+constants.Colon+dbDo.DataSetConfig.Port+
		dbDo.DataSetConfig.PreCreateDatasetContext, preDatasetBody, query, heads)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	preDatasetResp := mocks.CreatePreDatasetResp()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&preDatasetResp)
	if err != nil {
		return
	}

	return preDatasetResp, err
}
