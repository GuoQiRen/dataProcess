package clientApp

import (
	"dataProcess/constants"
	"dataProcess/constants/plat"
	"dataProcess/tools/utils/utils"
	"errors"
	"os"
	"runtime"
	"strings"
)

func (j *jinnFile) createSavePath(jinnPath, inOut, dataType string) (savePath string, err error) {
	savePath, err = j.getSaveFilePath(jinnPath, inOut, dataType)
	if err != nil {
		return
	}

	if j.inputMarkVersion == -1 && dataType == constants.MarkData { // 只创建前向目录
		shodSaveJsonPath := strings.Split(savePath, string(os.PathSeparator))
		shodSaveJsonPath = shodSaveJsonPath[:len(shodSaveJsonPath)-1]
		saveJsonPath := strings.Join(shodSaveJsonPath, string(os.PathSeparator))
		err = pathNotExistCreate(saveJsonPath)
	} else {
		err = pathNotExistCreate(savePath)
	}

	return
}

func (j *jinnFile) getSaveFilePath(jinnPath, inOut, dataType string) (savePath string, err error) {
	if len(jinnPath) < 1 {
		err = errors.New("jinnPath is too short")
		return
	}

	jinnPathArr := strings.Split(jinnPath, plat.LinuxSpiltRex)
	tailContext := jinnPathArr[len(jinnPathArr)-1]
	methodNameIndex := j.methodName + constants.SubLine + utils.IntToString(j.index)

	saveJsonPaths := make([]string, 1)
	if inOut == constants.InputData {
		if dataType == constants.MarkData {
			saveJsonPaths = append(saveJsonPaths, constants.Root, constants.Preprocessing,
				j.taskId, methodNameIndex, constants.InputData, constants.MarkData, constants.Resource, tailContext)
		}

		if dataType == constants.Material {
			saveJsonPaths = append(saveJsonPaths, constants.Root, constants.Preprocessing,
				j.taskId, methodNameIndex, constants.InputData, constants.Material, constants.Resource, tailContext)
		}
	}

	if inOut == constants.OutPutData {
		if dataType == constants.MarkData {
			saveJsonPaths = append(saveJsonPaths, constants.Root, constants.Preprocessing,
				j.taskId, methodNameIndex, constants.OutPutData, constants.MarkData, constants.Result, tailContext)
		}

		if dataType == constants.Material {
			saveJsonPaths = append(saveJsonPaths, constants.Root, constants.Preprocessing,
				j.taskId, methodNameIndex, constants.OutPutData, constants.Material, constants.Result, tailContext)
		}

	}

	switch runtime.GOOS {
	case plat.Windows:
		savePath = plat.EDisk + strings.Join(saveJsonPaths, plat.WindowsSpiltRex) // 这个是下载文件保存路径
	case plat.Linux:
		savePath = strings.Join(saveJsonPaths, plat.LinuxSpiltRex)
	}

	return
}

func (j *jinnFile) getJsonPath(curIndex int, inOut string) (jsonPath string, err error) {
	methodJsonPaths := make([]string, 1)
	methodNameIndex := j.methodName + constants.SubLine + utils.IntToString(curIndex)

	switch inOut {
	case constants.InputData:
		methodJsonPaths = append(methodJsonPaths, constants.Root, constants.Preprocessing, j.taskId,
			methodNameIndex, constants.InputData, constants.MarkData, constants.Resource)
	case constants.OutPutData:
		methodJsonPaths = append(methodJsonPaths, constants.Root, constants.Preprocessing, j.taskId,
			methodNameIndex, constants.OutPutData, constants.MarkData, constants.Result)
	}
	jsonPaths := methodJsonPaths[:len(methodJsonPaths)-2] // 真正保存methodJson的路径

	switch runtime.GOOS {
	case plat.Windows:
		jsonPath = plat.EDisk + strings.Join(jsonPaths, plat.WindowsSpiltRex) // 这个是方法连接json路径
	case plat.Linux:
		jsonPath = strings.Join(jsonPaths, plat.LinuxSpiltRex)
	}
	return
}

func (j *jinnFile) getMaterialAndMarkDataPath(materialMarkData, inOut string, downStage bool) (materialMarkDataPath string) {
	var methodNameIndex string
	materialMarkDataPaths := make([]string, 1)

	if downStage {
		methodNameIndex = j.methodName + constants.SubLine + utils.IntToString(j.index)

		if materialMarkData == constants.Material && inOut == constants.InputData {

			materialInPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.InputData, constants.Material, constants.Resource)
			materialMarkDataPath = strings.Join(materialInPath, string(os.PathSeparator))
		}

		if materialMarkData == constants.Material && inOut == constants.OutPutData {
			materialOutPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.OutPutData, constants.Material, constants.Result)
			materialMarkDataPath = strings.Join(materialOutPath, string(os.PathSeparator))
		}

		if materialMarkData == constants.MarkData && inOut == constants.InputData {
			markDataInPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.InputData, constants.MarkData, constants.Resource)
			materialMarkDataPath = strings.Join(markDataInPath, string(os.PathSeparator))
		}

		if materialMarkData == constants.MarkData && inOut == constants.OutPutData {
			markDataOutPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.OutPutData, constants.MarkData, constants.Result)
			materialMarkDataPath = strings.Join(markDataOutPath, string(os.PathSeparator))
		}

	} else {
		methodNameIndex = j.nextMethodName + constants.SubLine + utils.IntToString(j.index+1)

		if materialMarkData == constants.Material && inOut == constants.OutPutData {
			materialOutPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.OutPutData, constants.Material, constants.Result)
			materialMarkDataPath = strings.Join(materialOutPath, string(os.PathSeparator))
		}

		if materialMarkData == constants.MarkData && inOut == constants.OutPutData {
			markDataOutPath := append(materialMarkDataPaths, constants.Root, constants.Preprocessing, j.taskId,
				methodNameIndex, constants.OutPutData, constants.MarkData, constants.Result)
			materialMarkDataPath = strings.Join(markDataOutPath, string(os.PathSeparator))
		}
	}

	if runtime.GOOS == plat.Windows {
		materialMarkDataPath = plat.EDisk + materialMarkDataPath
	}

	return
}

func pathNotExistCreate(path string) (err error) {
	back, err := utils.PathExists(path)
	if err != nil {
		return
	}

	if !back {
		switch runtime.GOOS {
		case plat.Windows:
			err = os.MkdirAll(path, os.ModeDir)
		case plat.Linux:
			err = RecycleMkdir(path)
		}
		if err != nil {
			return
		}
	}
	return
}

func RecycleMkdir(path string) (err error) {
	dirs := strings.Split(path, plat.LinuxSpiltRex)
	length := len(dirs)
	initPath := plat.LinuxSpiltRex
	for i := 0; i < length; i++ {
		if dirs[i] != "" {
			initPath += dirs[i] + plat.LinuxSpiltRex
		}
		if initPath == plat.LinuxSpiltRex {
			continue
		}
		back, _ := utils.PathExists(initPath)
		if back {
			continue
		}
		err = os.MkdirAll(initPath, os.ModeDir)
		if err != nil {
			return
		}
		err = os.Chmod(initPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	return
}
