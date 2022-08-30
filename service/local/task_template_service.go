package local

import "dataProcess/model"

type TemplateManage interface {
	SaveTemplate(template model.Template) (id int32, err error)
	UpdateTemplate(template model.Template) (err error)
	SelectTemplateByUserId(userId int32) (templates []model.Template, err error)
	SelectTemplate(id int32) (template model.Template, err error)
	DeleteTemplate(id int32, userId int32) (err error)
	SelectAllTemplate(userId int32) (templates []model.Template, err error)
}
