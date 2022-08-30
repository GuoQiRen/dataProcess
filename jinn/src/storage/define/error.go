package define

const (
	// 资源
	ECPermissionDenied int32 = 1
	ECSpaceNotExist    int32 = 2
	ECPathInvalid      int32 = 3
	ECDirExist         int32 = 4
	ECDirNotExist      int32 = 5
	ECFileExist        int32 = 6
	ECFileNotExist     int32 = 7
	ECHandleInvalid    int32 = 8
	ECEndOfFile        int32 = 9
	// 参数
	ECParamError      int32 = 101
	ECDstNotDir       int32 = 102
	ECDstSameSrc      int32 = 103
	ECDstIsSubDir     int32 = 104
	ECFileNameInvalid int32 = 105
	// 其他
	ECDatabaseError   int32 = 201
	ECFileServerError int32 = 202
	ECUnknown         int32 = 9999
	// 网络协议
	ECSendError       int32 = 10000
	ECJsonFormatError int32 = 10001
	ECUnknownMethod   int32 = 10002
)

var codeMsg = map[int32]string{
	ECPermissionDenied: "permission denied",
	ECSpaceNotExist:    "space not exist",
	ECPathInvalid:      "path invalid",
	ECDirExist:         "dir exist",
	ECDirNotExist:      "dir not exist",
	ECFileExist:        "filej exist",
	ECFileNotExist:     "filej not exist",
	ECHandleInvalid:    "handle invalid",
	ECEndOfFile:        "EOF",
	ECParamError:       "",
	ECDstNotDir:        "dst not dir",
	ECDstSameSrc:       "dst same src",
	ECDstIsSubDir:      "dst is sub dir",
	ECFileNameInvalid:  "filej name irregular",
	ECDatabaseError:    "",
	ECFileServerError:  "",
	ECUnknown:          "unknown",
	ECSendError:        "",
	ECJsonFormatError:  "JSON format is incorrect",
	ECUnknownMethod:    "",
}

type Error = *errorImp

func NewError(code int32) Error {
	return &errorImp{Code: code}
}

type errorImp struct {
	Code   int32  `json:"code"`
	ErrMsg string `json:"errMsg"`
}

func (e *errorImp) Error() string {
	if len(e.ErrMsg) == 0 {
		return codeMsg[e.Code]
	}
	return e.ErrMsg
}

func (e *errorImp) SetInfo(s string) *errorImp {
	e.ErrMsg = s
	return e
}
