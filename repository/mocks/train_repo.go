package mocks

type TokenInfo struct {
	Token string `json:"token"`
}

type LoginResp struct {
	Result    bool      `json:"result"`
	ErrorMsg  string    `json:"errorMsg"`
	ErrorCode uint32    `json:"errorCode"`
	Data      TokenInfo `json:"data"`
}

func CreateLoginResp() LoginResp {
	return LoginResp{}
}

type TrainResp struct {
	Result           bool   `json:"result"`
	Id               uint32 `json:"id"`
	MissionCreatedId string `json:"missionCreatedId"`
	ErrorMsg         string `json:"errorMsg"`
}

func CreateJobResp() TrainResp {
	return TrainResp{}
}

type HostPathMap struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
}

type Env struct {
	EnvKey   string      `json:"envKey"`
	EnvValue interface{} `json:"envValue"`
}

type ContainerInfo struct {
	ImageInfo     string        `json:"imageInfo"`
	Hostname      string        `json:"hostname"`
	ContainerIp   string        `json:"containerIp"`
	ContainerName string        `json:"containerName"`
	Resource      interface{}   `json:"resource"`
	HostPathMaps  []HostPathMap `json:"hostPathMaps"`
	Envs          []Env         `json:"env"`
}

type RuntimeInfo struct {
	NodeInfo       []interface{}   `json:"nodeInfo"`
	ContainerInfos []ContainerInfo `json:"containerInfo"`
}

type RuntimeResp struct {
	Result      bool        `json:"result"`
	RuntimeInfo RuntimeInfo `json:"runtimeInfo"`
}

func CreateRuntimeResp() RuntimeResp {
	return RuntimeResp{}
}

type DataSetResp struct {
	Result    bool   `json:"result"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode uint32 `json:"errorCode"`
	VerDir    string `json:"verDir"`
	RealId    int    `json:"real_id"`
}

func CreateDataSetResp() DataSetResp {
	return DataSetResp{}
}

type TaskOperateResp struct {
	Result    bool   `json:"result"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode uint32 `json:"errorCode"`
}

func CreateTaskOperateResp() TaskOperateResp {
	return TaskOperateResp{}
}

type TaskTrainLogResp struct {
	Result      bool   `json:"result"`
	ErrorMsg    string `json:"errorMsg"`
	ErrorCode   string `json:"errorCode"`
	FileContent string `json:"fileContent"`
}

func CreateTaskTrainLogResp() TaskTrainLogResp {
	return TaskTrainLogResp{}
}

type MemberInfo struct {
	UserAccount string `json:"userAccount"`
	UserName    string `json:"userName"`
	MemberRole  int32  `json:"memberRole"`
}

type UserGroupContext struct {
	Result    bool         `json:"result"`
	ErrorMsg  string       `json:"errorMsg"`
	ErrorCode int32        `json:"errorCode"`
	Data      []MemberInfo `json:"data"`
}

func CreateUserGroupResp() UserGroupContext {
	return UserGroupContext{}
}

type GroupInfo struct {
	Id         int32  `json:"id"`
	MemberRole int32  `json:"memberRole"`
	GroupName  string `json:"groupName"`
}

type RoleInfo struct {
	RoleId          int32  `json:"roleId"`
	RoleName        string `json:"roleName"`
	PlatformId      int32  `json:"platformId"`
	RoleTag         string `json:"roleTag"`
	RoleDescription string `json:"roleDescription"`
	ParentRoleId    int32  `json:"parentRoleId"`
	Creator         string `json:"creator"`
	CreateTime      string `json:"createTime"`
	RoleStatus      int32  `json:"roleStatus"`
}

type RoleGroupInfo struct {
	GroupInfo []GroupInfo `json:"groupInfo"`
	RoleInfo  []RoleInfo  `json:"roleInfo"`
}

type UserDetailInfo struct {
	Id          int32         `json:"id"`
	UserAccount string        `json:"userAccount"`
	UserName    string        `json:"userName"`
	UserDomain  string        `json:"userDomain"`
	UserEmail   string        `json:"userEmail"`
	Creator     string        `json:"creator"`
	CreateTime  string        `json:"createTime"`
	UserStatus  int32         `json:"userStatus"`
	Data        RoleGroupInfo `json:"data"`
}

type UserDetailContext struct {
	Result    bool           `json:"result"`
	ErrorMsg  string         `json:"errorMsg"`
	ErrorCode int32          `json:"errorCode"`
	Data      UserDetailInfo `json:"data"`
}

func CreateUserDetailResp() UserDetailContext {
	return UserDetailContext{}
}

type UserInfo struct {
	UserAccount string
	UserDomain  string
	Exp         int64
	UserName    string
}

type TokenUpdateContext struct {
	Token     string   `json:"token"`
	ErrorCode int32    `json:"errorCode"`
	ErrorMsg  string   `json:"errorMsg"`
	UserInfo  UserInfo `json:"user_info"`
	Result    bool     `json:"result"`
}

func CreateUserTokenResp() TokenUpdateContext {
	return TokenUpdateContext{}
}
