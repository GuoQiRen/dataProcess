package clientApp

import (
	"dataProcess/constants"
	"dataProcess/constants/plat"
	"dataProcess/repository/mocks"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func (j *jinnFile) reCycleGenerateMethodJson(preMethodName, nextMethodName string) (jsonInPath string, err error) {
	var nextOutPath string
	nextLocalOutPaths := make([]string, 0)
	curIndex := j.index

	for _, path := range j.inputPaths {
		nextOutPath, err = j.createSavePath(path.InPath, constants.OutPutData, constants.Material)
		if err != nil {
			return
		}

		nextLocalOutPaths = append(nextLocalOutPaths, nextOutPath)
	}

	curMethodName := j.methodName
	jsonOutPath, err := j.getJsonPath(curIndex, constants.OutPutData)
	if err != nil {
		return
	}

	if j.index == 0 { // 表示第一次输入
		j.methodName = curMethodName
		jsonInPath, err = j.getJsonPath(curIndex, constants.InputData)
		if err != nil {
			return
		}
	} else { // 表示后续的输出
		j.methodName = preMethodName
		jsonInPath, err = j.getJsonPath(curIndex-1, constants.OutPutData)
		if err != nil {
			return
		}
	}

	statBody, err := readMethodJson(casePlatFilePath(jsonInPath, j.taskId+constants.JsonSuffix))
	if err != nil {
		return
	}

	var splitFileArr, splitJsonArr []string
	var curStatFiles []mocks.StatisticFile
	var curStatBody mocks.StatisticBody

	for _, statFile := range statBody.Files {
		// 如果下载存在.json结尾的，理论上认为是标注数据，不给予写入
		if statFile.FileOutPath[len(statFile.FileOutPath)-5:len(statFile.FileOutPath)] == constants.JsonSuffix {
			continue
		}

		switch runtime.GOOS {
		case plat.Windows:
			splitFileArr = strings.Split(statFile.FileOutPath, plat.WindowsSpiltRex) // 获取后缀文件名
			splitJsonArr = strings.Split(statFile.JsonOutPath, plat.WindowsSpiltRex)
		case plat.Linux:
			splitFileArr = strings.Split(statFile.FileOutPath, plat.LinuxSpiltRex) // 获取后缀文件名
			splitJsonArr = strings.Split(statFile.JsonOutPath, plat.LinuxSpiltRex)
		}

		curStatFile := mocks.StatisticFile{
			BaseId:   statFile.BaseId,
			ParentId: statFile.ParentId,
			FileId:   statFile.FileId,
			File:     statFile.FileOutPath,
			Json:     statFile.JsonOutPath,
		}

		// 第4个index索引下标要替换成output_data
		if len(splitJsonArr) != 0 && len(splitJsonArr) > 5 {
			splitJsonArr[4] = nextMethodName + constants.SubLine + utils.IntToString(curIndex+1)
			splitJsonArr[5] = constants.OutPutData

			jsonShouPath := strings.Join(splitJsonArr[:len(splitJsonArr)-1], string(os.PathSeparator))

			// 创建不存在的路径
			if curIndex+1 < j.total {
				err = pathNotExistCreate(jsonShouPath)
				if err != nil {
					return
				}
			}

			curStatFile.JsonOutPath = strings.Join(splitJsonArr, string(os.PathSeparator))
		}

		if len(splitFileArr) != 0 && len(splitFileArr) > 5 {
			splitFileArr[4] = nextMethodName + constants.SubLine + utils.IntToString(curIndex+1)
			splitFileArr[5] = constants.OutPutData

			fileShouPath := strings.Join(splitFileArr[:len(splitFileArr)-1], string(os.PathSeparator))

			// 创建不存在的路径
			if curIndex+1 < j.total {
				err = pathNotExistCreate(fileShouPath)
				if err != nil {
					return
				}
			}

			curStatFile.FileOutPath = strings.Join(splitFileArr, string(os.PathSeparator))
		}

		curStatFiles = append(curStatFiles, curStatFile)
	}

	curStatBody.Files = curStatFiles
	curStatBody.Args = j.taskArgs
	curStatBody.InputDirIds = statBody.InputDirIds
	curStatBody.MaterialInPath = statBody.MaterialOutPath
	curStatBody.MaterialOutPath = j.getMaterialAndMarkDataPath(constants.Material, constants.OutPutData, false)
	curStatBody.MarkDataInPath = statBody.MarkDataOutPath
	curStatBody.MarkDataOutPath = j.getMaterialAndMarkDataPath(constants.MarkData, constants.OutPutData, false)
	curStatBody.InputSpaceId = statBody.InputSpaceId
	curStatBody.OutputSpaceId = statBody.OutputSpaceId
	curStatBody.OutputDirId = statBody.OutputDirId

	if !j.inputMark { // 不需要下载数据，直接写在jsonFile里面
		curStatBody.InputSources = j.inputSources()
	}

	if !j.outputMark { // 不需要上传数据，直接写在jsonFile里面
		curStatBody.OutputSource = j.outputSources()
	}

	err = utils.GenerateJsonFile(curStatBody, casePlatFilePath(jsonOutPath, j.taskId+constants.JsonSuffix))
	if err != nil {
		return
	}

	return
}

func readMethodJson(jsonFile string) (statBody mocks.StatisticBody, err error) {
	jsonF, err := os.Open(jsonFile)
	if err != nil {
		return
	}
	defer jsonF.Close()

	jsonData, err := ioutil.ReadAll(jsonF)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &statBody)
	if err != nil {
		return
	}
	return
}
