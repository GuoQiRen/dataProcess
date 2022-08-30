package define

const (
	TabIdGen         string = "id_generator"
	TabDevice        string = "storage_device"
	TabFileBaseInfo  string = "file_base_info"
	TabSubFileInfo   string = "sub_file_info"
	TabUserSpaceList string = "user_space_list"
	TabUserSpace     string = "user_space_%d"
	TabFileStats     string = "file_stats_%d"
	TabLabelData     string = "label_data_%d"
	TabTransferLog   string = "transfer_log"
)

const (
	IdGenTable  string = "_id"
	IdGenNextId string = "next_id"
)

const (
	DeviceHost string = "_id"
)

const (
	FileBaseId         string = "_id"
	FileBaseDeviceId   string = "did"
	FileBaseParentId   string = "pid"
	FileBaseClass      string = "cl"
	FileBaseExtName    string = "en"
	FileBaseSize       string = "si"
	FileBaseSha1       string = "sh"
	FileBaseUpdateTime string = "ut"
	FileBaseWidth      string = "w"
	FileBaseHeight     string = "h"
	FileBaseFrameTotal string = "ft"
	FileBaseFrameRate  string = "fr"
	FileBaseDuration   string = "du"
)

const (
	FileClassUnknown string = "unknown"
	FileClassImage   string = "image"
	FileClassAudio   string = "audio"
	FileClassVideo   string = "video"
	FileClassOther   string = "other"
	FileClassThumb   string = "thumb"
	FileClassPreview string = "preview"
	FileClassFrame   string = "frame"
)

const (
	FileClassIDUnknown int32 = 0
	FileClassIDImage   int32 = 1
	FileClassIDAudio   int32 = 2
	FileClassIDVideo   int32 = 3
	FileClassIDOther   int32 = 9
	FileClassIDThumb   int32 = 101
	FileClassIDPreview int32 = 102
	FileClassIDFrame   int32 = 103
)

const (
	SubFileBaseId     string = "bi"
	SubFileDeviceId   string = "di"
	SubFileParentId   string = "pid"
	SubFileClass      string = "cl"
	SubFileExtName    string = "en"
	SubFileSize       string = "si"
	SubFileUpdateTime string = "ut"
	SubFileWidth      string = "w"
	SubFileHeight     string = "h"
	SubFileFrameNum   string = "fn"
)

const (
	UserSpaceId         string = "_id"
	UserSpaceName       string = "name"
	UserSpaceLabel      string = "label"
	UserSpaceHide       string = "hide"
	UserSpaceNextFileId string = "next_id"
	UserSpaceStatus     string = "statuses"
)

const (
	FileInfoId         string = "_id"
	FileInfoDir        string = "d"
	FileInfoType       string = "t"
	FileInfoName       string = "n"
	FileInfoUserPerm   string = "pu"
	FileInfoGroupPerm  string = "pg"
	FileInfoAdmin      string = "pa"
	FileInfoPreview    string = "pv"
	FileInfoRead       string = "pr"
	FileInfoAppend     string = "pn"
	FileInfoModify     string = "pm"
	FileInfoDelete     string = "pd"
	FileInfoCreator    string = "cr"
	FileInfoCreateTime string = "ct"
	FileInfoUpdateTime string = "ut"
	FileInfoBaseId     string = "bi"
	FileInfoDeviceId   string = "di"
	FileInfoClass      string = "cl"
	FileInfoExtName    string = "en"
	FileInfoSize       string = "si"
	FileInfoWidth      string = "w"
	FileInfoHeight     string = "h"
	FileInfoFrameTotal string = "ft"
	FileInfoFrameRate  string = "fr"
	FileInfoDuration   string = "du"
)

const (
	FileTypeUnknown int32 = 0
	FileTypeDir     int32 = 1
	FileTypeFile    int32 = 2
)

const (
	LabelTagKey   string = "k"
	LabelTagValue string = "v"
	LabelTagGroup string = "g"
)

const (
	ObjDataId     string = "id"
	ObjDataParent string = "pa"
	ObjDataClass  string = "cl"
	ObjDataShape  string = "sh"
	ObjDataCoord  string = "co"
	ObjDataProps  string = "pr"
)

const (
	LabelVer        string = "ver"
	LabelBaseId     string = "id"
	LabelFrameNum   string = "fn"
	LabelImgSize    string = "si"
	LabelTags       string = "ta"
	LabelObjects    string = "ob"
	LabelNextId     string = "ni"
	LabelUpdateTime string = "ut"
)

func GetClassId(cls string) int32 {
	switch cls {
	case FileClassUnknown:
		return FileClassIDUnknown
	case FileClassImage:
		return FileClassIDImage
	case FileClassAudio:
		return FileClassIDAudio
	case FileClassVideo:
		return FileClassIDVideo
	case FileClassThumb:
		return FileClassIDThumb
	case FileClassPreview:
		return FileClassIDPreview
	case FileClassFrame:
		return FileClassIDFrame
	default:
		return FileClassIDOther
	}
}

func GetClassString(cls int32) string {
	switch cls {
	case FileClassIDUnknown:
		return FileClassUnknown
	case FileClassIDImage:
		return FileClassImage
	case FileClassIDAudio:
		return FileClassAudio
	case FileClassIDVideo:
		return FileClassVideo
	case FileClassIDThumb:
		return FileClassThumb
	case FileClassIDPreview:
		return FileClassPreview
	case FileClassIDFrame:
		return FileClassFrame
	default:
		return FileClassOther
	}
}

func ConvertFileType(t int32, d int) int32 {
	if d > 0 {
		switch t {
		case FileTypeDir:
			return 0
		case FileTypeFile:
			return 1
		default:
			return t
		}
	} else {
		switch t {
		case 0:
			return FileTypeDir
		case 1:
			return FileTypeFile
		default:
			return t
		}
	}
}
