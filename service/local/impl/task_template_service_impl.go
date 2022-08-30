package impl

import (
	"dataProcess/constants/atoi"
	"dataProcess/constants/table"
	"dataProcess/model"
	"dataProcess/orm"
	dbDo "dataProcess/tcp_socket/db_do"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TaskTemplateImpl struct {
	coll *mgo.Collection
}

func CreateTaskTemplateImpl() *TaskTemplateImpl {
	if orm.MongoDb == nil {
		mongoDb := new(dbDo.MongodbServer)
		orm.MongoDb = mongoDb.SetUp()
	}
	return &TaskTemplateImpl{coll: orm.MongoDb.C(table.TaskTemplate)}
}

func TemplateCompletePostDecode(c *gin.Context) (templateBody model.Template, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&templateBody)
	return templateBody, err
}

func (t *TaskTemplateImpl) SaveTemplate(template model.Template) (id int32, err error) {
	taskTemplateProxy := CreateAtomicIdProxy(nil, t, nil, nil)
	id, err = taskTemplateProxy.GetAtomicId(atoi.Template)
	if err != nil {
		return
	}

	template.Id = id
	template.CreateTime = utils.TimeToString(time.Now())
	template.UpdateTime = utils.TimeToString(time.Now())
	err = t.coll.Insert(template)
	if err != nil {
		return
	}
	return
}

func (t *TaskTemplateImpl) UpdateTemplate(template model.Template) (err error) {
	err = t.coll.Update(&bson.M{"id": template.Id},
		bson.M{"$set": bson.M{
			"name":                     template.Name,
			"inputType":                template.InputType,
			"inputPath":                template.InputPath,
			"outputType":               template.OutputType,
			"outputPath":               template.OutputPath,
			"arguments":                template.Arguments,
			"updateTime":               utils.TimeToString(time.Now()),
			"description":              template.Description,
			"inputMarkVersion":         template.InputMarkVersion,
			"outputDataSetVersionName": template.OutputDataSetVersionName,
			"outputMarkVersion":        template.OutputMarkVersion,
		}})
	if err != nil {
		return
	}
	return
}

func (t *TaskTemplateImpl) SelectTemplate(id int32) (template model.Template, err error) {
	err = t.coll.Find(bson.M{"id": id}).One(&template)
	return
}

func (t *TaskTemplateImpl) SelectTemplateByUserId(userId int32) (templates []model.Template, err error) {
	err = t.coll.Find(bson.M{"createId": userId}).All(&templates)
	return
}

func (t *TaskTemplateImpl) DeleteTemplate(id, userId int32) (err error) {
	tempDetail, err := t.SelectTemplate(id)
	if err != nil {
		return
	}
	if tempDetail.CreateId != userId {
		err = errors.New("permission denied")
		return
	}
	err = t.coll.Remove(bson.M{"id": id})
	if err != nil {
		return
	}
	return
}

func (t *TaskTemplateImpl) SelectAllTemplate(userId int32) (templates []model.Template, err error) {
	templates, err = t.SelectTemplateByUserId(userId)
	return
}
