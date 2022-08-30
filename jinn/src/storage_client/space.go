package client

import (
	"dataProcess/jinn/src/storage/database"
	"dataProcess/jinn/src/utils/jsonx"
)

const (
	MethodGetSpaceId     = "get_space_id"
	MethodGetPathInfo    = "get_path_info"
	MethodGetDirStatInfo = "get_dir_stat_info"
	MethodList           = "list"
	MethodMkdir          = "mkdir"
	MethodGetPerm        = "get_perm"
	MethodCopySelect     = "copy_select"
	MethodCopyDir        = "copy_dir"
	MethodUploadCreate   = "upload_create"
	MethodUploadWrite    = "upload_write"
	MethodDirOpen        = "dir_open"
	MethodDirRead        = "dir_read"
	MethodDirClose       = "dir_close"
	MethodFileOpen       = "file_open"
	MethodFileGetInfo    = "file_get_info"
	MethodFileSeek       = "file_seek"
	MethodFileRead       = "file_read"
	MethodFileClose      = "file_close"
	MethodFinderCreate   = "finder_create"
	MethodFinderRead     = "finder_read"
	MethodFinderClose    = "finder_close"
)

type GetSpaceIdParams struct {
	Name string
}

func (p *GetSpaceIdParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "name", p.Name)
	return jsonx.AppendDocumentEnd(doc)
}

type GetSpaceIdResult struct {
	SpaceId int64
}

func (r *GetSpaceIdResult) Decode(cur *jsonx.Cursor) {
	r.SpaceId = cur.Int64()
}

type GetPathInfoParams struct {
	Space int64
	Path  string
}

func (p *GetPathInfoParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendStringElement(doc, "path", p.Path)
	return jsonx.AppendDocumentEnd(doc)
}

type GetPathInfoResult struct {
	database.FileInfo
}

func (r *GetPathInfoResult) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		key := cur.Key()
		switch key {
		case "id":
			r.Id = cur.Int64()
		case "dir":
			r.Dir = cur.Int64()
		case "type":
			r.Type = cur.Int32()
		case "name":
			r.Name = cur.String()
		case "creator":
			r.Creator = cur.Int32()
		}
	}
}

type MkdirParams struct {
	Space int64
	Dir   int64
	Name  string
}

func (p *MkdirParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64Element(doc, "dir", p.Dir)
	doc = jsonx.AppendStringElement(doc, "name", p.Name)
	return jsonx.AppendDocumentEnd(doc)
}

type MkdirResult struct {
	DirId int64
}

func (r *MkdirResult) Decode(cur *jsonx.Cursor) {
	r.DirId = cur.Int64()
}

type GetPermParams struct {
	Space int64
	Dir   int64
}

func (p *GetPermParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64Element(doc, "dir", p.Dir)
	return jsonx.AppendDocumentEnd(doc)
}

type GetPermResult = PassThroughResult

type CopySelectParams struct {
	DstSpace int64
	DstDir   int64
	DstToken string
	SrcSpace int64
	SrcFile  []int64
	Option   struct {
		Cover bool
		Rule  string
		Start int32
	}
}

func (p *CopySelectParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "dstSpace", p.DstSpace)
	doc = jsonx.AppendInt64Element(doc, "dstDir", p.DstDir)
	doc = jsonx.AppendStringElement(doc, "dstToken", p.DstToken)
	doc = jsonx.AppendInt64Element(doc, "srcSpace", p.SrcSpace)
	doc = jsonx.AppendInt64ArrayElement(doc, "srcFile", p.SrcFile)
	doc = jsonx.AppendDocumentElementStart(doc, "option")
	doc = jsonx.AppendBoolElement(doc, "cover", p.Option.Cover)
	doc = jsonx.AppendStringElement(doc, "rule", p.Option.Rule)
	doc = jsonx.AppendInt32Element(doc, "start", p.Option.Start)
	doc = jsonx.AppendDocumentEnd(doc)
	return jsonx.AppendDocumentEnd(doc)
}

type CopyDirParams struct {
	DstSpace int64
	DstDir   int64
	DstToken string
	SrcSpace int64
	SrcDir   int64
	Option   struct {
		Filter string
		Cover  bool
		Rule   string
		Start  int32
	}
}

func (p *CopyDirParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "dstSpace", p.DstSpace)
	doc = jsonx.AppendInt64Element(doc, "dstDir", p.DstDir)
	doc = jsonx.AppendStringElement(doc, "dstToken", p.DstToken)
	doc = jsonx.AppendInt64Element(doc, "srcSpace", p.SrcSpace)
	doc = jsonx.AppendInt64Element(doc, "srcDir", p.SrcDir)
	doc = jsonx.AppendDocumentElementStart(doc, "option")
	doc = jsonx.AppendStringElement(doc, "filter", p.Option.Filter)
	doc = jsonx.AppendBoolElement(doc, "cover", p.Option.Cover)
	doc = jsonx.AppendStringElement(doc, "rule", p.Option.Rule)
	doc = jsonx.AppendInt32Element(doc, "start", p.Option.Start)
	doc = jsonx.AppendDocumentEnd(doc)
	return jsonx.AppendDocumentEnd(doc)
}

type FileOperateResult struct {
	Ok   int32
	Skip int32
}

func (r *FileOperateResult) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "ok":
			r.Ok = cur.Int32()
		case "skip":
			r.Skip = cur.Int32()
		}
	}
}

type UploadCreateParams struct {
	Space int64
	Dir   int64
	Name  string
	Size  int64
}

func (p *UploadCreateParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64Element(doc, "dir", p.Dir)
	doc = jsonx.AppendStringElement(doc, "name", p.Name)
	doc = jsonx.AppendInt64Element(doc, "size", p.Size)
	return jsonx.AppendDocumentEnd(doc)
}

type UploadCreateResult struct {
	Handle uint32
	BaseId int64
}

func (r *UploadCreateResult) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "handle":
			r.Handle = cur.Uint32()
		case "baseId":
			r.BaseId = cur.Int64()
		}
	}
}

type DirOpenParams struct {
	Space int64
	Dir   int64
}

func (p *DirOpenParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64Element(doc, "dir", p.Dir)
	return jsonx.AppendDocumentEnd(doc)
}

type FileOpenParams struct {
	Space    int64
	Dir      int64
	Name     string
	SubType  int32
	FrameNum int32
}

func (p *FileOpenParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64Element(doc, "dir", p.Dir)
	doc = jsonx.AppendStringElement(doc, "name", p.Name)
	doc = jsonx.AppendInt32Element(doc, "subType", p.SubType)
	doc = jsonx.AppendInt32Element(doc, "frameNum", p.FrameNum)
	return jsonx.AppendDocumentEnd(doc)
}

type FileGetInfoResult struct {
	Info database.FileInfo
}

func (r *FileGetInfoResult) Decode(cur *jsonx.Cursor) {
	r.Info.Decode(cur)
}

type FileSeekParams struct {
	Handle uint32
	Offset int64
	Whence int
}

func (p *FileSeekParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "handle", p.Handle)
	doc = jsonx.AppendInt64Element(doc, "offset", p.Offset)
	doc = jsonx.AppendInt32Element(doc, "whence", int32(p.Whence))
	return jsonx.AppendDocumentEnd(doc)
}

type FileSeekResult struct {
	Offset int64
}

func (r *FileSeekResult) Decode(cur *jsonx.Cursor) {
	r.Offset = cur.Int64()
}
