package mocks

type StatisticFile struct {
	BaseId      int64  `json:"baseId"`
	FileId      int64  `json:"fileId"`
	ParentId    int64  `json:"parentId"`
	File        string `json:"file"`
	FileOutPath string `json:"fileOutPath"`
	Json        string `json:"json"`
	JsonOutPath string `json:"jsonOutPath"`
}

type StatisticBody struct {
	InputSpaceId    int64                  `json:"inputSpaceId"`
	InputDirIds     []int64                `json:"inputDirIds"`
	OutputSpaceId   int64                  `json:"outputSpaceId"`
	OutputDirId     int64                  `json:"outputDirId"`
	MaterialInPath  string                 `json:"materialInPath"`
	MaterialOutPath string                 `json:"materialOutPath"`
	MarkDataInPath  string                 `json:"markDataInPath"`
	MarkDataOutPath string                 `json:"markDataOutPath"`
	Files           []StatisticFile        `json:"files"`
	Args            map[string]interface{} `json:"args"`
	InputSources    []string               `json:"inputSources"`
	OutputSource    string                 `json:"outputSources"`
}
