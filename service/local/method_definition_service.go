package local

import (
	"dataProcess/model"
)

type MethodDefinition interface {
	UploadDefinitionMethod(upMethodOrm model.UPMethodOrm) (id int32, err error)
	UploadDefinitionMethodUpdate(methodId, userId int32, upMethodOrm model.UPMethodOrm) (err error)
	AssignDefinitionMethodAuthority(methodId, userId string, authInfos model.AuthInfo) (err error)
	SelectDefinitionMethodAuthority(methodId, userId string) (authInfos model.AuthInfo, err error)
	DeleteDefinitionMethod(methodId, userId string) (err error)
}
