package method

/**
算子类型分为两种，一种是公共算子，另一种是我的算子，共享给我的算子为3, 公共算子为2，我的算子为1
*/
const (
	MyMethod = iota + 1
	ShareMethod
	PubMethod
	ALLMethod = 777
)
