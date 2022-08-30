package impl

import (
	"dataProcess/constants"
	"dataProcess/constants/atoi"
	"dataProcess/constants/filej"
	"dataProcess/constants/method"
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
	"strings"
)

type CategoryImpl struct {
	coll *mgo.Collection
}

func CategoryDecode(c *gin.Context) (cateOrm model.CategoryOrm, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&cateOrm)
	return cateOrm, err
}

func GetPostFormValues(c *gin.Context, jsonTags []string) (values []string) {
	values = make([]string, 0, len(jsonTags))
	for _, jsTag := range jsonTags {
		values = append(values, c.DefaultPostForm(jsTag, ""))
	}
	return
}

func CategoryIdsDecode(c *gin.Context) (ids []int32, err error) {
	idds := struct {
		Ids []int32 `json:"ids" bson:"ids"`
	}{}
	err = json.NewDecoder(c.Request.Body).Decode(&idds)
	return idds.Ids, err
}

func CreateCategoryImpl() *CategoryImpl {
	if orm.MongoDb == nil {
		mongoDb := new(dbDo.MongodbServer)
		orm.MongoDb = mongoDb.SetUp()
	}
	return &CategoryImpl{coll: orm.MongoDb.C(table.MethodCategory)}
}

func (c *CategoryImpl) CreateCategory(categoryOrm model.CategoryOrm) (err error) {
	atomicIdProxy := CreateAtomicIdProxy(nil, nil, c, nil)

	cateId, err := atomicIdProxy.GetAtomicId(atoi.Method)
	if err != nil {
		return
	}

	categoryOrm.Id = cateId
	err = c.coll.Insert(&categoryOrm)
	if err != nil {
		return
	}

	return
}

func (c *CategoryImpl) UpdateCategory(categoryOrm model.CategoryOrm) (err error) {
	curCateOrm, err := c.SelectCategoryById(categoryOrm.Id)
	if err != nil {
		return
	}

	super, err := CertainUserSuperManage(utils.IntToString(int(categoryOrm.CreatorId)))
	if err != nil {
		return
	}

	if categoryOrm.CreatorId != curCateOrm.CreatorId && !super {
		err = errors.New("you have not permission modify it")
		return
	}

	err = c.coll.Update(bson.M{"id": categoryOrm.Id}, bson.M{"$set": bson.M{"name": categoryOrm.Name}})
	if err != nil {
		return
	}
	return
}

func (c *CategoryImpl) SelectCategoryById(id int32) (categoryOrm model.CategoryOrm, err error) {
	err = c.coll.Find(bson.M{"id": id}).One(&categoryOrm)
	return
}

func (c *CategoryImpl) SelectCategoryByIds(ids []int32) (categoryOrms []model.CategoryOrm, err error) {
	if len(ids) == 0 { // 查询所有
		err = c.coll.Find(nil).All(&categoryOrms)
		return
	}

	var query []bson.M
	for _, id := range ids {
		query = append(query, bson.M{"id": id})
	}
	err = c.coll.Find(bson.M{"$or": query}).All(&categoryOrms)
	return
}

func (c *CategoryImpl) SelectCategoryByParentId(parentId int32) (categoryOrms []model.CategoryOrm, err error) {
	err = c.coll.Find(bson.M{"parentId": parentId}).All(&categoryOrms)
	return
}

func (c *CategoryImpl) SelectCategoryByName(name string) (categoryOrms []model.CategoryOrm, err error) {
	query := bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: name, Options: "i"}}}
	err = c.coll.Find(query).All(&categoryOrms)
	return
}

func (c *CategoryImpl) SelectCategoryIds() (ids []int32, err error) {
	var categoryOrms []model.CategoryOrm
	err = c.coll.Find(nil).All(&categoryOrms)
	for _, categoryOrm := range categoryOrms {
		ids = append(ids, categoryOrm.Id)
	}
	return
}

func (c *CategoryImpl) SelectAllCategory(root *model.CategoryTree, name, userId string) (matchTree []model.CategoryTree, err error) {
	userIdInt32, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	super, err := CertainUserSuperManage(userId)
	if err != nil {
		return
	}

	if len(name) == 0 {
		QueryAllCategories(root, root.Id, userIdInt32, super)
		root.Children = CountSort(root)
		ShareChildrenAssign(root) // 查询分享给我的算子
		SuperChildrenAssign(root) // 查询分享给我的算子

		if !super { // 非超级管理员只显示前三项
			root.Children = root.Children[:len(root.Children)-1]
		}
		return
	}

	matchTree, err = QueryMatchCategories(name, userIdInt32, super)
	return
}

func (c *CategoryImpl) DeleteCategory(id, userId int32) (err error) {
	catOrm, err := c.SelectCategoryById(id)
	if err != nil {
		return
	}

	cateNext, err := c.SelectCategoryByParentId(catOrm.Id)
	if err != nil {
		return
	}

	if len(cateNext) != 0 {
		err = errors.New("you can not delete this node with sub node")
		return
	}

	super, err := CertainUserSuperManage(utils.IntToString(int(userId)))
	if err != nil {
		return
	}

	if catOrm.Level == filej.MethodType {
		auth := CreateDefinitionImpl().AuthOrNotAccess(id, userId, method.Deleted) ||
			CreateDefinitionImpl().AuthOrNotAccess(id, userId, method.Completed)

		if !auth && !super {
			err = errors.New("this method: " + utils.IntToString(int(id)) + " is permission denied")
			return
		}
	} else {
		if userId != catOrm.CreatorId && !super {
			err = errors.New("you have not permission delete it")
			return
		}
	}

	err = c.coll.Remove(bson.M{"id": id})
	if err != nil {
		return
	}
	return
}

func QueryMatchCategories(name string, userId int32, super bool) (matchTree []model.CategoryTree, err error) {
	matchCategories, err := CreateCategoryImpl().SelectCategoryByName(name)
	if err != nil {
		return
	}

	authBd := DefinitionImpl{MethodImpl: CreateMethodImpl()}
	for _, matchCategory := range matchCategories {
		cur := model.CategoryTree{Id: matchCategory.Id, ParentId: matchCategory.ParentId,
			Name: matchCategory.Name, Children: make([]*model.CategoryTree, 0)}
		methodInfo, _ := CreateMethodImpl().SelectMethodById(matchCategory.Id)
		if methodInfo.DeleteMark {
			continue
		}

		if !super { // 并非管理员查看，有限制
			if methodInfo.Type == method.ShareMethod || methodInfo.Type == method.MyMethod { // 是分类
				// 删除的不能搜索到，还有是我没有权限的不能搜索到
				if !authBd.AuthOrNotAccess(methodInfo.Id, userId, method.Used) &&
					!authBd.AuthOrNotAccess(methodInfo.Id, userId, method.Deleted) &&
					!authBd.AuthOrNotAccess(methodInfo.Id, userId, method.Completed) { // 共享给我的
					continue
				}
			}
		}

		QueryAllCategories(&cur, matchCategory.Id, userId, super) // 从父id往下递归查询
		matchTree = append(matchTree, cur)
	}
	return
}

func QueryAllCategories(root *model.CategoryTree, rootIndex, userId int32, super bool) {
	myChildren := make([]*model.CategoryTree, 0)
	shareChildren := make([]*model.CategoryTree, 0)
	shareRemainCategoryTree := new([]*model.CategoryTree)
	superChildren := make([]*model.CategoryTree, 0)
	publicChildren := make([]*model.CategoryTree, 0)
	markRemainIds := make([]int32, 0)

	QueryMyCategories(&myChildren, method.MyMethod, userId)
	QueryPublicCategories(&publicChildren, method.PubMethod, userId)
	QueryShareCategories(&shareChildren, &markRemainIds, method.MyMethod, userId)
	FilterUseLessCategory(markRemainIds, &shareChildren, shareRemainCategoryTree)
	QuerySuperCategories(&superChildren, method.MyMethod, userId, super)

	LinkCategories(root, rootIndex, myChildren, *shareRemainCategoryTree, superChildren, publicChildren)
}

func LinkCategories(root *model.CategoryTree, rootIndex int32, myChildren, shareChildren,
	superChildren, publicChildren []*model.CategoryTree) {

	children, err := CreateCategoryImpl().SelectCategoryByParentId(rootIndex)
	if err != nil {
		return
	}
	for _, child := range children {
		newNode := model.CategoryTree{Id: child.Id, Name: child.Name, ParentId: child.ParentId,
			CreatorId: child.CreatorId, CreatorName: child.CreatorName, Level: child.Level,
			Children: make([]*model.CategoryTree, 0)}
		switch child.Id {
		case method.MyMethod:
			newNode.Children = myChildren
		case method.ShareMethod:
			newNode.Children = shareChildren
		case method.PubMethod:
			newNode.Children = publicChildren
		case method.ALLMethod:
			newNode.Children = superChildren
		}
		root.Children = append(root.Children, &newNode)
	}
}

func QueryPublicCategories(publicChildren *[]*model.CategoryTree, rootIndex, userId int32) {
	children, err := CreateCategoryImpl().SelectCategoryByParentId(rootIndex)
	if err != nil {
		return
	}

	if len(children) == 0 {
		return
	}

	for _, child := range children {
		newNode := model.CategoryTree{Id: child.Id, Name: child.Name, ParentId: child.ParentId,
			CreatorId: child.CreatorId, CreatorName: child.CreatorName, Level: child.Level,
			Children: make([]*model.CategoryTree, 0)}
		*publicChildren = append(*publicChildren, &newNode)
	}

	// 我上传的算子
	for _, child := range *publicChildren {
		QueryPublicCategories(&child.Children, child.Id, userId)
	}
}

func QueryShareCategories(shareChildren *[]*model.CategoryTree, markRemainIds *[]int32, rootIndex, userId int32) {
	children, err := CreateCategoryImpl().SelectCategoryByParentId(rootIndex)
	if err != nil {
		return
	}

	if len(children) == 0 {
		return
	}

	methodProxy := CreateMethodImpl()
	authBd := DefinitionImpl{MethodImpl: methodProxy}

	for _, child := range children {
		if child.CreatorId != userId { // 并未我上传的
			if child.Level == filej.MethodType {
				if authBd.AuthOrNotAccess(child.Id, userId, method.Used) ||
					authBd.AuthOrNotAccess(child.Id, userId, method.Deleted) ||
					authBd.AuthOrNotAccess(child.Id, userId, method.Completed) { // 共享给我的,并且我有权限
					*markRemainIds = append(*markRemainIds, child.Id)
				}
			}

			newNode := model.CategoryTree{Id: child.Id, Name: child.Name, ParentId: child.ParentId,
				CreatorId: child.CreatorId, CreatorName: child.CreatorName, Level: child.Level,
				Children: make([]*model.CategoryTree, 0)}
			*shareChildren = append(*shareChildren, &newNode)
		}

	}

	// 我上传的算子
	for _, child := range *shareChildren {
		QueryShareCategories(&child.Children, markRemainIds, child.Id, userId)
	}
}

func FilterUseLessCategory(remainIds []int32, shareChildren, remainCategoryTree *[]*model.CategoryTree) {
	SaveShareMethod(remainIds, shareChildren, remainCategoryTree)
}

func SaveShareMethod(remainIds []int32, shareChildren, remainCategoryTree *[]*model.CategoryTree) {
	for _, child := range *shareChildren {
		if child.Level == filej.MethodType { // 如果是方法
			ok := utils.Int32ContainKeyOrNot(child.Id, remainIds)
			if ok { // 删除
				*remainCategoryTree = append(*remainCategoryTree, child)
			}
		} else {
			SaveShareMethod(remainIds, &child.Children, remainCategoryTree)
		}
	}

	return
}

func DeleteEmptyCate(shareChildren *[]*model.CategoryTree) {
	for index, child := range *shareChildren {
		if child.Level == filej.CateType && len(child.Children) == 0 { // 如果下面没有算子了，并且当前为分类，直接删除
			if index+1 >= len(*shareChildren) {
				*shareChildren = (*shareChildren)[:len(*shareChildren)-1]
			} else if index == 0 {
				*shareChildren = (*shareChildren)[index+1:]
			} else {
				*shareChildren = append((*shareChildren)[:index], (*shareChildren)[index+1:]...)
			}
		}
		DeleteEmptyCate(&child.Children)
	}
}

func QueryMyCategories(myChildren *[]*model.CategoryTree, rootIndex, userId int32) {
	children, err := CreateCategoryImpl().SelectCategoryByParentId(rootIndex)
	if err != nil {
		return
	}

	if len(children) == 0 {
		return
	}

	var methodInfo model.MethodOrm

	methodProxy := CreateMethodImpl()

	for _, child := range children {
		if child.Level == filej.MethodType {
			methodInfo, err = methodProxy.SelectMethodById(child.Id)
			if err != nil {
				return
			}
			if methodInfo.DeleteMark {
				continue
			}
		}
		if child.CreatorId == userId { // 我上传的
			newNode := model.CategoryTree{Id: child.Id, Name: child.Name, ParentId: child.ParentId,
				CreatorId: child.CreatorId, CreatorName: child.CreatorName, Level: child.Level,
				Children: make([]*model.CategoryTree, 0)}
			*myChildren = append(*myChildren, &newNode)
		}
	}

	// 我上传的算子
	for _, child := range *myChildren {
		QueryMyCategories(&child.Children, child.Id, userId)
	}
}

func QuerySuperCategories(superChildren *[]*model.CategoryTree, rootIndex, userId int32, super bool) {
	children, err := CreateCategoryImpl().SelectCategoryByParentId(rootIndex)
	if err != nil {
		return
	}

	if len(children) == 0 {
		return
	}

	var methodInfo model.MethodOrm

	methodProxy := CreateMethodImpl()
	for _, child := range children {
		if child.Level == filej.MethodType {
			methodInfo, err = methodProxy.SelectMethodById(child.Id)
			if err != nil {
				return
			}
			if methodInfo.DeleteMark {
				continue
			}
		}

		if super {
			newSuperNode := model.CategoryTree{Id: child.Id, Name: child.Name, ParentId: child.ParentId,
				CreatorId: child.CreatorId, CreatorName: child.CreatorName, Level: child.Level,
				Children: make([]*model.CategoryTree, 0)}
			*superChildren = append(*superChildren, &newSuperNode)
		}
	}

	for _, child := range *superChildren {
		QuerySuperCategories(&child.Children, child.Id, userId, super)
	}
}

/**
按索引排序
*/
func CountSort(root *model.CategoryTree) (sortRootChildren []*model.CategoryTree) {
	sortRootChildren = make([]*model.CategoryTree, len(root.Children))

	for _, child := range root.Children {
		if int(child.Id)-1 >= len(sortRootChildren) {
			sortRootChildren[len(sortRootChildren)-1] = child
		} else {
			sortRootChildren[child.Id-1] = child
		}
	}
	return
}

func ShareChildrenAssign(root *model.CategoryTree) {
	shareTree := root.Children[method.ShareMethod-1]

	newInsertUserIdCategoryNode(shareTree)

	for _, child := range root.Children {
		if child.Id == method.ShareMethod && shareTree != nil {
			child.Children = shareTree.Children
			break
		}
	}
}

func SuperChildrenAssign(root *model.CategoryTree) {
	superTree := root.Children[len(root.Children)-1]

	// 插入id_name的分类节点
	newInsertUserIdCategoryNode(superTree)

	for _, child := range root.Children {
		if child.Id == method.ALLMethod && superTree != nil {
			child.Children = superTree.Children
			break
		}
	}
}

func newInsertUserIdCategoryNode(nodeTree *model.CategoryTree) {

	nodeMap := make(map[string][]*model.CategoryTree)

	// 多加一层userId+_+userName分类, superTree 按照 creatorId 分组
	for _, nodeChild := range nodeTree.Children {
		nodeKey := utils.IntToString(int(nodeChild.CreatorId)) + constants.SubLine + nodeChild.CreatorName
		if nodeChildren, ok := nodeMap[nodeKey]; ok {
			nodeChildren = append(nodeChildren, nodeChild)
			nodeMap[nodeKey] = nodeChildren
		} else {
			newChildren := make([]*model.CategoryTree, 1)
			newChildren[0] = nodeChild
			nodeMap[nodeKey] = newChildren
		}
	}

	// 写回superTree
	nodeTree.Children = make([]*model.CategoryTree, 0)

	for superKey, categoryTrees := range nodeMap {
		superKeys := strings.Split(superKey, constants.SubLine)
		createIdInt, err := utils.StringToint32(superKeys[0])
		if err != nil {
			return
		}

		newCateNode := model.CategoryTree{Id: createIdInt, Name: superKey, ParentId: nodeTree.Id,
			CreatorId: createIdInt, CreatorName: superKeys[1], Level: filej.CateType,
			Children: make([]*model.CategoryTree, 0)}
		for index, cateTree := range categoryTrees {
			cateTree.ParentId = createIdInt
			categoryTrees[index] = cateTree
		}
		newCateNode.Children = append(newCateNode.Children, categoryTrees...)
		nodeTree.Children = append(nodeTree.Children, &newCateNode)
	}
}
