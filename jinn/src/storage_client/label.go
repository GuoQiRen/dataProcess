package client

import (
	"dataProcess/jinn/src/storage/database"
	"dataProcess/jinn/src/utils/jsonx"
)

const (
	MethodGetTags      = "label.get_tags"
	MethodGetFileObjs  = "label.get_file_objs"
	MethodGetFrameObjs = "label.get_frame_objs"
)

type GetTagsParams struct {
	Ver   int64
	Space int64
	File  []int64
}

func (p *GetTagsParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64ArrayElement(doc, "filej", p.File)
	return jsonx.AppendDocumentEnd(doc)
}

type GetFileObjsParams struct {
	Ver    int64
	BaseId []int64
	Filter *jsonx.Cursor
}

func (p *GetFileObjsParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendInt64ArrayElement(doc, "baseId", p.BaseId)
	return jsonx.AppendDocumentEnd(doc)
}

type GetFrameObjsParams struct {
	Ver        int64
	BaseId     int64
	FrameRange [2]int32
	Filter     *jsonx.Cursor
}

func (p *GetFrameObjsParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendInt64Element(doc, "baseId", p.BaseId)
	doc = jsonx.AppendInt32ArrayElement(doc, "frameRange", p.FrameRange[:])
	return jsonx.AppendDocumentEnd(doc)
}

type LabelDataResult struct {
	Infos []database.LabelInfo
}

func (r *LabelDataResult) Decode(cur *jsonx.Cursor) {
	r.Infos = make([]database.LabelInfo, cur.Size())
	for i := 0; cur.Next(); i++ {
		r.Infos[i].Decode(cur.Value())
	}
}
