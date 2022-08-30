package model

type InputPaths struct {
	InPath     string `json:"inPath" bson:"inPath"`
	SeqId      string `json:"seqId" bson:"seqId"`
	SolidId    int32  `json:"solidId" bson:"solidId"`
	PrefixPath string `json:"prefixPath" bson:"prefixPath"`
}

type OutputPaths struct {
	OutPath    string `json:"outPath" bson:"outPath"`
	SeqId      string `json:"seqId" bson:"seqId"`
	SolidId    int32  `json:"solidId" bson:"solidId"`
	PrefixPath string `json:"prefixPath" bson:"prefixPath"`
}

type arguments struct {
	MethodId   int32                  `json:"methodId" bson:"methodId"`
	MethodName string                 `json:"methodName" bson:"methodName"`
	Args       map[string]interface{} `json:"args" bson:"args"`
}

type TaskOrm struct {
	Name                     string        `json:"name" bson:"name"`
	InputType                int32         `json:"inputType" bson:"inputType"`
	InputPath                []InputPaths  `json:"inputPath" bson:"inputPath"`
	OutputType               int32         `json:"outputType" bson:"outputType"`
	OutputPath               []OutputPaths `json:"outputPath" bson:"outputPath"`
	InputMarkVersion         int64         `json:"inputMarkVersion" bson:"inputMarkVersion"`
	OutputDataSetVersionName string        `json:"outputDataSetVersionName" bson:"outputDataSetVersionName"`
	OutputMarkVersion        int64         `json:"outputMarkVersion" bson:"outputMarkVersion"`
	DownloadContent          []string      `json:"downloadContent" bson:"downloadContent"`
	MaterialType             []string      `json:"materialType" bson:"materialType"`
	SrcDirName               string        `json:"srcDirName" bson:"srcDirName"`
	Arguments                []arguments   `json:"arguments" bson:"arguments"`
	TaskProtoTypeId          int32         `json:"taskProtoTypeId" bson:"taskProtoTypeId"`
	CreateId                 int32         `json:"createId" bson:"createId"`
	CreateName               string        `json:"createName" bson:"createName"`
	Description              string        `json:"description" bson:"description"`
}

type TaskSaveOrm struct {
	Id                       int32         `json:"id" bson:"id"`
	UId                      string        `json:"uId"`
	Name                     string        `json:"name" bson:"name"`
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

func CreateTaskSaveOrmTrain(id int32, trainRelatedId string) TaskSaveOrm {
	return TaskSaveOrm{Id: id, TrainRelatedId: trainRelatedId}
}

func CreateTaskSaveOrmRecall(id int32, status, stage, exception, endTime string) TaskSaveOrm {
	return TaskSaveOrm{Id: id, Status: status, Stage: stage, Exception: exception, EndTime: endTime}
}

func CreateTaskSaveOrmReplyDataSet(id int32, inputPath []InputPaths) TaskSaveOrm {
	return TaskSaveOrm{Id: id, InputPath: inputPath}
}

func CreateTaskInputMarkVersion(id int32, inputMarkVersion int64) TaskSaveOrm {
	return TaskSaveOrm{Id: id, InputMarkVersion: inputMarkVersion}
}
