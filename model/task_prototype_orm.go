package model

type Proto struct {
	Links         interface{} `json:"links"`
	OperatorTypes interface{} `json:"operatorTypes"`
	Operators     interface{} `json:"operators"`
}

type ProtoType struct {
	Id    int32 `json:"id"`
	Proto Proto `json:"proto"`
}
