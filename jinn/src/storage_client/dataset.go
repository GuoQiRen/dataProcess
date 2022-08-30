package client

import (
	"dataProcess/jinn/src/utils/jsonx"
)

const (
	MethodDatasetCreate     = "dataset.create"
	MethodDatasetAppendFile = "dataset.append_file"
	MethodDatasetDirOpen    = "dataset.dir_open"
	MethodDatasetWriteLabel = "dataset.write_label"
	MethodDatasetFindLabel  = "dataset.find_label"
	MethodDatasetLabelOpen  = "dataset.label_open"
	MethodDatasetLabelRead  = "dataset.label_read"
	MethodDatasetLabelClose = "dataset.label_close"
)

type DatasetFindLabelParams struct {
	Ver    int64
	Space  int64
	File   []int64
	Filter *jsonx.Cursor
	Fields *jsonx.Cursor
}

func (p *DatasetFindLabelParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendInt64Element(doc, "space", p.Space)
	doc = jsonx.AppendInt64ArrayElement(doc, "filej", p.File)
	doc = append(doc, `"fields":["ver","frameNo","imgSize","tags","objects"]`...)
	return jsonx.AppendDocumentEnd(doc)
}

type DatasetLabelOpenParams struct {
	Ver  int64
	Kind string
}

func (p *DatasetLabelOpenParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "ver", p.Ver)
	doc = jsonx.AppendStringElement(doc, "kind", p.Kind)
	return jsonx.AppendDocumentEnd(doc)
}
