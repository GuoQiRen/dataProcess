package impl

import (
	"dataProcess/constants/atoi"
	"dataProcess/constants/table"
	"dataProcess/model"
	"dataProcess/orm"
	dbDo "dataProcess/tcp_socket/db_do"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TaskPrototypeImpl struct {
	coll *mgo.Collection
}

func PrototypeDecode(c *gin.Context) (protoType model.ProtoType, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&protoType)
	return protoType, err
}

func CreateTaskPrototypeImpl() *TaskPrototypeImpl {
	if orm.MongoDb == nil {
		mongoDb := new(dbDo.MongodbServer)
		orm.MongoDb = mongoDb.SetUp()
	}
	return &TaskPrototypeImpl{coll: orm.MongoDb.C(table.TaskPrototype)}
}

func (t *TaskPrototypeImpl) SaveProto(protoType model.ProtoType) (protoTypeRet model.ProtoType, err error) {
	protoTypeAtomic := CreateAtomicIdProxy(nil, nil, nil, t)
	protoId, err := protoTypeAtomic.GetAtomicId(atoi.Proto)
	if err != nil {
		return
	}
	protoType.Id = protoId
	err = t.coll.Insert(protoType)
	return protoType, err
}

func (t *TaskPrototypeImpl) UpdateProto(protoType model.ProtoType) (err error) {
	err = t.coll.Update(&bson.M{"id": protoType.Id}, bson.M{"$set": bson.M{"proto": protoType.Proto}})
	if err != nil {
		return
	}
	return
}

func (t *TaskPrototypeImpl) SelectProto(id int32) (protoType model.ProtoType, err error) {
	err = t.coll.Find(bson.M{"id": id}).One(&protoType)
	return
}

func (t *TaskPrototypeImpl) DeleteProto(id int32) (err error) {
	err = t.coll.Remove(bson.M{"id": id})
	if err != nil {
		return
	}
	return
}
