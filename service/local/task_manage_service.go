package local

import (
	"dataProcess/model"
	"dataProcess/repository/mocks"
)

type TaskManage interface {
	CopyTask(id, userId, userName string) (newId int32, err error)
	TerminalTask(id string) (err error)
	DeleteRemoteTask(id string) (err error)
	StartTask(id string, userId string) (err error)
	DeleteTask(id, userId int32) (err error)
	LogTask(id string) (taskLog string, err error)
	SaveTask(taskOrm model.TaskOrm) (id int32, err error)
	UpdateTask(id int32, taskOrm model.TaskOrm) (err error)
	UpdateTaskStatus(taskSaveOrm model.TaskSaveOrm) (err error)
	UpdateTaskTrainId(taskSaveOrm model.TaskSaveOrm) (err error)
	UpdateTaskInputPathSeqId(taskSaveOrm model.TaskSaveOrm) (err error)
	UpdateTaskInputMarkVersionSeqId(taskSaveOrm model.TaskSaveOrm) (err error)
	QueryTaskDetail(id int32) (taskSaveOrm model.TaskSaveOrm, err error)
	SearchTask(queryTask *mocks.QueryTask) (taskSaveOrms []model.TaskSaveOrm, err error)
}
