package model

import "mime/multipart"

type UPMethodOrm struct {
	EnName             string         `json:"enName" bson:"enName"`
	CnName             string         `json:"cnName" bson:"cnName"`
	ExtName            string         `json:"extName" bson:"extName"`
	CpuNums            int32          `json:"cpuNums" bson:"cpuNums"`
	GpuNums            int32          `json:"gpuNums" bson:"gpuNums"`
	GpuType            string         `json:"gpuType" bson:"gpuType"`
	Arguments          []interface{}  `json:"arguments" bson:"arguments"`
	InputMark          bool           `json:"inputMark" bson:"inputMark"`
	OutputMark         bool           `json:"outputMark" bson:"outputMark"`
	DealMarkData       bool           `json:"dealMarkData" bson:"dealMarkData"`
	DealMaterialType   string         `json:"dealMaterialType" bson:"dealMaterialType"`
	FileSelf           multipart.File `json:"fileSelf" bson:"fileSelf"`
	CreateId           int32          `json:"createId" bson:"createId"`
	CreateName         string         `json:"createName" bson:"createName"`
	CateParentId       int32          `json:"cateParentId" bson:"cateParentId"`
	UpDescription      multipart.File `json:"description" bson:"description"`
	DescriptionEnName  string         `json:"descriptionEnName" bson:"descriptionEnName"`
	DescriptionExtName string         `json:"descriptionExtName" bson:"descriptionExtName"`
}
