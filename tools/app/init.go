package app

type ReturnData struct {
	Result    bool        `json:"result"`
	ErrorMsg  string      `json:"errorMsg"`
	ErrorCode int32       `json:"errorCode"`
	Data      interface{} `json:"data"`
}

type PageResponse struct {
	ErrorCode int32       `json:"errorCode" example:"0"`
	Data      interface{} `json:"data"`
	ErrorMsg  string      `json:"errorMsg"`
	Result    bool        `json:"result"`
	Total     int32       `json:"total"`
}

func (r *ReturnData) CreateReturn(result bool, errMsg string, code int32) *ReturnData {
	if errMsg != "" {
		r.ErrorMsg = errMsg
	}
	r.ErrorCode = code
	r.Result = result
	return r
}

func (p *PageResponse) CreatePageResponse(result bool, errMsg string, code int32, total int32) *PageResponse {
	if errMsg != "" {
		p.ErrorMsg = errMsg
	}
	p.ErrorCode = code
	p.Result = result
	p.Total = total
	return p
}
