package define

const KeyLock string = "temp:lock:%s"

// space:id:dir:
const (
	SpaceCacheInfo string = "space:%d:%d:info"
	SpaceDirChange string = "space:dir_change"
)

const (
	StreamStat    string = "stream:stat"
	StreamPreproc string = "stream:preprocess"
	StreamFrame   string = "stream:frame"
)

const (
	TaskFrame string = "task:extract_frame:%d"
)
