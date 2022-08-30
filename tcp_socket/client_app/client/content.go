package clientApp

import (
	"dataProcess/constants/mark"
	"dataProcess/constants/plat"
	rpcCli "dataProcess/jinn/src/cgo_rpc_client"
	"dataProcess/jinn/src/storage/database"
	"dataProcess/tools/utils/utils"
	"strings"
)

func getSpaceIdInfo(jinnPath string, client rpcCli.Client) (spaceId int64, pathInfo *database.FileInfo, err error) {
	// 分割用户空间名称
	spaceName, pathContext, err := utils.SplitJinnString(jinnPath)
	if err != nil {
		return
	}

	// 获取用户空间id
	spaceId, err = client.GetSpaceId(spaceName)
	if err != nil {
		return
	}

	// 获取用户空间路径信息
	pathInfo, err = client.GetPathInfo(spaceId, pathContext)
	if err != nil {
		return
	}

	return
}

/**
获取统一存储spaceInfo
*/
func (j *jinnFile) getUniqueSpaceInfo(pathDir string) (spaceId int64, dirId int64, err error) {

	var pathInfo *database.FileInfo

	spaceId, pathInfo, err = getSpaceIdInfo(pathDir, client)
	if err != nil {
		return
	}

	// 数据集需要重置spaceId和dirId
	if pathInfo != nil {
		dirId = pathInfo.Id
	}

	return
}

/**
获取数据集的SpaceInfo
*/
func (j *jinnFile) getDataSetSpaceInfo(path, seqId, markVersion string) (spaceId int64, dirId int64, err error) {
	spaceId = 900
	var pathContext, realId string

	// 获取后缀context
	paths := strings.Split(path, plat.LinuxSpiltRex)
	if markVersion == mark.InputMarkVersion {
		inputMarkVersion := int(j.inputMarkVersion)
		if inputMarkVersion == 0 { // 寻找seqId作为realId
			realId = seqId
		} else {
			realId = utils.IntToString(int(j.inputMarkVersion))
		}
	}

	if markVersion == mark.OutputMarkVersion {
		outputMarkVersion := int(j.outputMarkVersion)
		if outputMarkVersion == 0 { // 寻找seqId作为realId
			realId = seqId
		} else {
			realId = utils.IntToString(int(j.outputMarkVersion))
		}
	}

	if len(paths) > 3 {
		paths = paths[3:]
		pathContext = plat.LinuxSpiltRex + realId + plat.LinuxSpiltRex + strings.Join(paths, plat.LinuxSpiltRex)
	} else {
		pathContext = plat.LinuxSpiltRex + realId
	}

	pathInfo, err := client.GetPathInfo(spaceId, pathContext) // 需要改成任一路径
	if err != nil {
		return
	}

	// 数据集需要重置spaceId和dirId
	if pathInfo != nil {
		dirId = pathInfo.Id
	}

	return
}
