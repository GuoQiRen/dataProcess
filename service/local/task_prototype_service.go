package local

import "dataProcess/model"

type TaskProtoType interface {
	SaveProto(protoType model.ProtoType) (protoTypeRet model.ProtoType, err error)
	UpdateProto(protoType model.ProtoType) (err error)
	SelectProto(id int32) (protoType model.ProtoType, err error)
	DeleteProto(id int32) (err error)
}
