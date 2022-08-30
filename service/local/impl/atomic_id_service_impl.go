package impl

import (
	"dataProcess/constants/atoi"
	"dataProcess/model"
	"errors"
	"sync"
)

var mutex sync.Mutex

type AtomicIdProxy struct {
	taskManageAtomicProxy     *TaskManageImpl
	taskTemplateAtomicProxy   *TaskTemplateImpl
	taskProtoAtomicProxy      *TaskPrototypeImpl
	methodCategoryAtomicProxy *CategoryImpl
}

func CreateAtomicIdProxy(taskManageImpl *TaskManageImpl, taskTemplateImpl *TaskTemplateImpl, methodCategoryAtomicProxy *CategoryImpl, taskProtoAtomicProxy *TaskPrototypeImpl) (atomicIdProxy AtomicIdProxy) {
	atomicIdProxy = AtomicIdProxy{
		taskManageAtomicProxy:     taskManageImpl,
		taskTemplateAtomicProxy:   taskTemplateImpl,
		methodCategoryAtomicProxy: methodCategoryAtomicProxy,
		taskProtoAtomicProxy:      taskProtoAtomicProxy}
	return
}

func (a *AtomicIdProxy) GetAtomicId(atomType int) (id int32, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	var taskOrmSaves []model.TaskSaveOrm
	var methodCategories []model.CategoryOrm
	var taskProtoTypes []model.ProtoType

	switch atomType {
	case atoi.Task:
		err = a.taskManageAtomicProxy.coll.Find(nil).Sort("id").All(&taskOrmSaves)
	case atoi.Template:
		err = a.taskTemplateAtomicProxy.coll.Find(nil).Sort("id").All(&taskOrmSaves)
	case atoi.Method:
		err = a.methodCategoryAtomicProxy.coll.Find(nil).Sort("id").All(&methodCategories)
	case atoi.Proto:
		err = a.taskProtoAtomicProxy.coll.Find(nil).Sort("id").All(&taskProtoTypes)
	default:
		err = errors.New("atomic id type not certain")
	}

	if err != nil {
		return
	}

	if len(taskOrmSaves) != 0 {
		id = taskOrmSaves[len(taskOrmSaves)-1].Id + 1
	} else if len(methodCategories) != 0 {
		id = methodCategories[len(methodCategories)-1].Id + 1
	} else if len(taskProtoTypes) != 0 {
		id = taskProtoTypes[len(taskProtoTypes)-1].Id + 1
	} else {
		id = 1
	}

	return
}
