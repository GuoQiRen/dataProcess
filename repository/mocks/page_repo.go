package mocks

type PageCore struct {
	PageNum  int32 `json:"pageNum"`
	PageSize int32 `json:"pageSize"`
}

func CreatePageCore(pageNum, pageSize int32) PageCore {
	return PageCore{PageNum: pageNum, PageSize: pageSize}
}
