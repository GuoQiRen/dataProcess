package client

import (
	"dataProcess/jinn/src/storage/database"
	"dataProcess/jinn/src/utils/jsonx"
)

type RequestMsg struct {
	Method string
	Params jsonx.Marshaler
	User   int32
	Group  []int32
}

func (m *RequestMsg) Marshal(doc []byte, id uint32) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "id", id)
	doc = jsonx.AppendStringElement(doc, "method", m.Method)
	doc = jsonx.AppendHeader(doc, "params")
	doc = m.Params.MarshalJson(doc)
	if m.User > 0 {
		doc = jsonx.AppendInt32Element(doc, "user", m.User)
		doc = jsonx.AppendInt32ArrayElement(doc, "group", m.Group)
	}
	return jsonx.AppendDocumentEnd(doc)
}

type EmptyResult struct{}

func (r *EmptyResult) Decode(_ *jsonx.Cursor) {}

type PassThroughParams jsonx.Cursor

func (p *PassThroughParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	return jsonx.AppendValue(doc, (*jsonx.Cursor)(p).String())
}

type PassThroughResult struct {
	Cur *jsonx.Cursor
}

func (r *PassThroughResult) Decode(cur *jsonx.Cursor) {
	r.Cur = cur.Clone()
}

type HandleCreateResult struct {
	Handle uint32
}

func (r *HandleCreateResult) Decode(cur *jsonx.Cursor) {
	r.Handle = cur.Uint32()
}

type HandleOperateParams struct {
	Handle uint32
}

func (p *HandleOperateParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "handle", p.Handle)
	return jsonx.AppendDocumentEnd(doc)
}

type HandleReadParams struct {
	Handle uint32
	Size   int32
}

func (p *HandleReadParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "handle", p.Handle)
	doc = jsonx.AppendInt32Element(doc, "size", p.Size)
	return jsonx.AppendDocumentEnd(doc)
}

type FileInfoResult struct {
	Infos []database.FileInfo
}

func (r *FileInfoResult) Decode(cur *jsonx.Cursor) {
	r.Infos = make([]database.FileInfo, cur.Size())
	for i := 0; cur.Next(); i++ {
		r.Infos[i].Decode(cur.Value())
	}
}

type AuthParams struct {
	Key string
}

func (p *AuthParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "key", p.Key)
	return jsonx.AppendDocumentEnd(doc)
}

type DatasetCreateParams struct {
	Ver int64
}

func (p *DatasetCreateParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	return jsonx.AppendDocumentEnd(doc)
}

type DatasetWriteLabelParams struct {
	Ver   int64
	Infos []database.LabelInfo
}

func (p *DatasetWriteLabelParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendArrayElementStart(doc, "infos")
	for i := range p.Infos {
		doc = p.Infos[i].MarshalJson(doc)
	}
	doc = jsonx.AppendArrayEnd(doc)
	return jsonx.AppendDocumentEnd(doc)
}
