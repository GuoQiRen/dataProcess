package clientApp

import "dataProcess/model"

type jinnFile struct {
	userId                   int32
	taskId                   string
	index                    int
	total                    int
	methodName               string
	nextMethodName           string
	inputMarkVersion         int64
	outputMarkVersion        int64
	inputType                int32
	inputPaths               []model.InputPaths
	outputType               int32
	outputPath               []model.OutputPaths
	inputMark                bool
	outputMark               bool
	taskArgs                 map[string]interface{}
	seqId                    int32
	solidId                  int32
	outputDataSetVersionName string
	srcDirName               string
	downloadContent          []string
	materialType             []string
	existDealMarkData        bool
}

func CreateJinnBody(taskId, methodName, nextMethodName string, index, total int, userId int32, taskArgs map[string]interface{},
	inputMarkVersion, outputMarkVersion int64, inputType int32, inputPaths []model.InputPaths, outputType int32, outputPath []model.OutputPaths, outputDataSetVersionName string,
	srcDirName string, downloadContent, materialType []string, existDealMarkData, inputMark, outputMark bool) *jinnFile {
	return &jinnFile{taskId: taskId, methodName: methodName, nextMethodName: nextMethodName, index: index, total: total,
		userId: userId, taskArgs: taskArgs, inputMarkVersion: inputMarkVersion, outputMarkVersion: outputMarkVersion,
		inputType: inputType, inputPaths: inputPaths, outputType: outputType, outputPath: outputPath,
		outputDataSetVersionName: outputDataSetVersionName, srcDirName: srcDirName, downloadContent: downloadContent,
		materialType: materialType, existDealMarkData: existDealMarkData, inputMark: inputMark,
		outputMark: outputMark}
}
