package mocks

type QueryTask struct {
	Id        int32    `json:"id" bson:"id"`
	UserId    int32    `json:"userId" bson:"userId"`
	Creators  []int32  `json:"creators" bson:"creators"`
	Name      string   `json:"name" bson:"name"`
	Statuses  []string `json:"statuses" bson:"statuses"`
	BeginTime string   `json:"beginTime" bson:"beginTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	PageNum   int32    `json:"pageNum" bson:"pageNum"`
	PageSize  int32    `json:"pageSize" bson:"pageSize"`
}

func CreateQueryTask() QueryTask {
	return QueryTask{}
}
