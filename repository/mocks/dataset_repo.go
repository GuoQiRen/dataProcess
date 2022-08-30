package mocks

type VersionResp struct {
	VersionId int32 `json:"version_id"`
}

type PreDatasetResp struct {
	Result    bool        `json:"result"`
	ErrorMsg  string      `json:"errorMsg"`
	ErrorCode uint32      `json:"errorCode"`
	Data      VersionResp `json:"data"`
}

func CreatePreDatasetResp() PreDatasetResp {
	return PreDatasetResp{}
}

type FormalDatasetResp struct {
	Result    bool   `json:"result"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode uint32 `json:"errorCode"`
}

func CreateFormalDatasetResp() FormalDatasetResp {
	return FormalDatasetResp{}
}
