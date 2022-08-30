package model

type TemplateBody struct {
	TaskSaveOrm TaskSaveOrm `json:"taskSaveOrm"`
}

type Template struct {
	Id                       int32         `json:"id" bson:"id"`
	UId                      string        `json:"uId"`
	Name                     string        `json:"name" bson:"name"`
	TemplateName             string        `json:"templateName" bson:"templateName"`
	Status                   string        `json:"statuses" bson:"statuses"`
	Stage                    string        `json:"stage" bson:"stage"`
	InputType                int32         `json:"inputType" bson:"inputType"`
	InputPath                []InputPaths  `json:"inputPath" bson:"inputPath"`
	OutputType               int32         `json:"outputType" bson:"outputType"`
	OutputPath               []OutputPaths `json:"outputPath" bson:"outputPath"`
	InputMarkVersion         int64         `json:"inputMarkVersion" bson:"inputMarkVersion"`
	OutputDataSetVersionName string        `json:"outputDataSetVersionName" bson:"outputDataSetVersionName"`
	OutputMarkVersion        int64         `json:"outputMarkVersion" bson:"outputMarkVersion"`
	SrcDirName               string        `json:"srcDirName" bson:"srcDirName"`
	DownloadContent          []string      `json:"downloadContent" bson:"downloadContent"`
	MaterialType             []string      `json:"materialType" bson:"materialType"`
	Arguments                []arguments   `json:"arguments" bson:"arguments"`
	CreateId                 int32         `json:"createId" bson:"createId"`
	CreateName               string        `json:"createName" bson:"createName"`
	CreateTime               string        `json:"createTime" bson:"createTime"`
	UpdateTime               string        `json:"updateTime" bson:"updateTime"`
	EndTime                  string        `json:"endTime" bson:"endTime"`
	Description              string        `json:"description" bson:"description"`
	TrainRelatedId           string        `json:"trainRelatedId" bson:"trainRelatedId"`
	TaskProtoTypeId          int32         `json:"taskProtoTypeId" bson:"taskProtoTypeId"`
	Exception                string        `json:"exception" bson:"exception"`
}
