package define

type M = map[string]interface{}
type Tags = map[string]map[string][]string

const (
	StatusOk     string = "ok"
	StatusCreate string = "create"
	StatusDelete string = "delete"
)

type Perm uint32

const (
	PermAdmin   Perm = 1 << 31
	PermPreview Perm = 1 << 0
	PermRead    Perm = 1 << 1
	PermAppend  Perm = 1 << 2
	PermModify  Perm = 1 << 3
	PermDelete  Perm = 1 << 4
)

const (
	SpUserSys int32 = -1
)

const (
	RootDirId int64 = 1
)
