package model

type CategoryOrm struct {
	Id          int32  `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	ParentId    int32  `json:"parentId" bson:"parentId"`
	CreatorId   int32  `json:"creatorId" bson:"creatorId"`
	CreatorName string `json:"creatorName" bson:"creatorName"`
	Level       int32  `json:"level" bson:"level"`
}

type CategoryTree struct {
	Id          int32           `json:"id"`
	Name        string          `json:"name"`
	ParentId    int32           `json:"parentId"`
	CreatorId   int32           `json:"creatorId" bson:"creatorId"`
	CreatorName string          `json:"creatorName" bson:"creatorName"`
	Level       int32           `json:"level" bson:"level"`
	Children    []*CategoryTree `json:"children"`
}
