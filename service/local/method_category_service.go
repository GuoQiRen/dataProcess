package local

import (
	"dataProcess/model"
)

type MethodCategory interface {
	CreateCategory(categoryOrm model.CategoryOrm) (err error)
	UpdateCategory(categoryOrm model.CategoryOrm) (err error)
	SelectCategoryById(id int32) (categoryOrm model.CategoryOrm, err error)
	SelectCategoryByIds(ids []int32) (categoryOrms []model.CategoryOrm, err error)
	SelectCategoryIds() (ids []int32, err error)
	SelectCategoryByParentId(parentId int32) (categoryOrms []model.CategoryOrm, err error)
	SelectCategoryByName(name string) (categoryOrms []model.CategoryOrm, err error)
	SelectAllCategory(root *model.CategoryTree, name, userId string) (matchTree []model.CategoryTree, err error)
	DeleteCategory(id, userId int32) (err error)
}
