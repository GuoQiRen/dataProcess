package model

type UserAuth struct {
	UserId         int32   `json:"userId" bson:"userId"`
	ControlAspects []int32 `json:"controlAspects" bson:"controlAspects"`
}

type UserGroupsAuth struct {
	UserGroupId    int32   `json:"userGroupId" bson:"userGroupId"`
	ControlAspects []int32 `json:"controlAspects" bson:"controlAspects"`
}

type AuthInfo struct {
	UserAuth       []UserAuth       `json:"userAuth" bson:"userAuth"`
	UserGroupsAuth []UserGroupsAuth `json:"userGroupsAuth" bson:"userGroupsAuth"`
}

type AuthInfos struct {
	AuthInfo AuthInfo `json:"authInfo" bson:"authInfo"`
}

type MethodOrm struct {
	Id               int32         `json:"id" bson:"id"`
	EnName           string        `json:"enName" bson:"enName"`
	CnName           string        `json:"cnName" bson:"cnName"`
	Type             int32         `json:"type" bson:"type"`
	InputMark        bool          `json:"inputMark" bson:"inputMark"`
	OutputMark       bool          `json:"outputMark" bson:"outputMark"`
	DeleteMark       bool          `json:"deleteMark" bson:"deleteMark"`
	DealMarkData     bool          `json:"dealMarkData" bson:"dealMarkData"`
	DealMaterialType string        `json:"dealMaterialType" bson:"dealMaterialType"`
	CpuNums          int32         `json:"cpuNums" bson:"cpuNums"`
	GpuNums          int32         `json:"gpuNums" bson:"gpuNums"`
	GpuType          string        `json:"gpuType" bson:"gpuType"`
	Arguments        []interface{} `json:"arguments" bson:"arguments"`
	AuthInfo         AuthInfo      `json:"authInfo" bson:"authInfo"`
	CreateId         int32         `json:"createId" bson:"createId"`
	CreateName       string        `json:"createName" bson:"createName"`
	CreateTime       string        `json:"createTime" bson:"createTime"`
	UpdateTime       string        `json:"updateTime" bson:"updateTime"`
	Description      string        `json:"description" bson:"description"`
	MainEngine       string        `json:"mainEngine" bson:"mainEngine"`
}
