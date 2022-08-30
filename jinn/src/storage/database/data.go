package database

import (
	"dataProcess/jinn/src/storage/define"
	"dataProcess/jinn/src/utils/jsonx"
	"time"
)

type FileBaseInfo struct {
	Id          int64
	DevId       int32
	ParentId    int64
	FileClass   int32
	FileExtName string
	FileSize    int64
	Sha1        string
	updateTime  time.Time
	ImgWidth    int32
	ImgHeight   int32
	FrameTotal  int32
	FrameRate   int32
	Duration    int32
}

type SubFileInfo struct {
	id         int64
	BaseId     int64
	DevId      int32
	ParentId   int64
	Class      int32
	ExtName    string
	Size       int64
	updateTime time.Time
	Width      int32
	Height     int32
	FrameNum   int32
}

type UserSpaceInfo struct {
	Id     int64
	Name   string
	Label  string
	hide   bool
	nextId int64
	status string
}

type UserPerm struct {
	Admin   []int32
	Preview []int32
	Read    []int32
	Append  []int32
	Modify  []int32
	Delete  []int32
}

type GroupPerm struct {
	Preview []int32
	Read    []int32
	Append  []int32
	Modify  []int32
	Delete  []int32
}

type DirPerm struct {
	User  UserPerm
	Group GroupPerm
}

type DirInfo struct {
	Id         int64
	Dir        int64
	Type       int32
	Name       string
	Perm       *DirPerm
	Creator    int32
	CreateTime time.Time
	updateTime time.Time
}

type FileInfo struct {
	DirInfo
	BaseId     int64
	DevId      int32
	Class      string
	ExtName    string
	Size       int64
	Width      int32
	Height     int32
	FrameTotal int32
	FrameRate  int32
	Duration   int32
}

func (d *FileInfo) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt64Element(doc, "id", d.Id)
	doc = jsonx.AppendInt64Element(doc, "dir", d.Dir)
	doc = jsonx.AppendInt32Element(doc, "type", d.Type)
	doc = jsonx.AppendStringElement(doc, "name", d.Name)
	doc = jsonx.AppendInt32Element(doc, "creator", d.Creator)
	doc = jsonx.AppendTimestampElement(doc, "createTime", d.CreateTime)
	if d.Type == define.FileTypeFile {
		doc = jsonx.AppendInt64Element(doc, "baseId", d.BaseId)
		doc = jsonx.AppendStringElement(doc, "class", d.Class)
		doc = jsonx.AppendStringElement(doc, "extName", d.ExtName)
		doc = jsonx.AppendInt64Element(doc, "size", d.Size)
		doc = jsonx.AppendInt32Element(doc, "width", d.Width)
		doc = jsonx.AppendInt32Element(doc, "height", d.Height)
		doc = jsonx.AppendInt32Element(doc, "frameTotal", d.FrameTotal)
		doc = jsonx.AppendInt32Element(doc, "frameRate", d.FrameRate)
		doc = jsonx.AppendInt32Element(doc, "duration", d.Duration)
	}
	return jsonx.AppendDocumentEnd(doc)
}

func (d *FileInfo) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "id":
			d.Id = cur.Int64()
		case "dir":
			d.Dir = cur.Int64()
		case "type":
			d.Type = cur.Int32()
		case "name":
			d.Name = cur.String()
		case "creator":
			d.Creator = cur.Int32()
		case "createTime":
			d.CreateTime = time.Unix(cur.Int64(), 0)
		case "baseId":
			d.BaseId = cur.Int64()
		case "class":
			d.Class = cur.String()
		case "extName":
			d.ExtName = cur.String()
		case "size":
			d.Size = cur.Int64()
		case "width":
			d.Width = cur.Int32()
		case "height":
			d.Height = cur.Int32()
		case "frameTotal":
			d.FrameTotal = cur.Int32()
		case "frameRate":
			d.FrameRate = cur.Int32()
		case "duration":
			d.Duration = cur.Int32()
		}
	}
}

type ObjectInfo struct {
	Id     int32               `json:"id"`
	Parent []int32             `json:"parent"`
	Class  string              `json:"class"`
	Shape  string              `json:"shape"`
	Coord  [][2]int32          `json:"coord"`
	Props  map[string][]string `json:"props"`
}

type LabelInfo struct {
	Ver        string       `json:"ver"`
	BaseId     int64        `json:"baseId"`
	FileName   string       `json:"fileName"`
	SchemaVer  int32        `json:"schemaVer"`
	FrameNum   int32        `json:"frameNum"`
	ImgSize    [2]int32     `json:"imgSize"`
	Tags       define.Tags  `json:"tags"`
	Objects    []ObjectInfo `json:"objects"`
	nextId     int32
	updateTime time.Time
}

type LabelTags define.Tags

func (d *ObjectInfo) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendInt32Element(doc, "id", d.Id)
	doc = jsonx.AppendInt32ArrayElement(doc, "parent", d.Parent)
	doc = jsonx.AppendStringElement(doc, "class", d.Class)
	doc = jsonx.AppendStringElement(doc, "shape", d.Shape)
	doc = jsonx.AppendArrayElementStart(doc, "coord")
	for _, point := range d.Coord {
		doc = jsonx.AppendArrayStart(doc)
		doc = jsonx.AppendInt32(doc, int32(point[0]))
		doc = jsonx.AppendInt32(doc, int32(point[1]))
		doc = jsonx.AppendArrayEnd(doc)
	}
	doc = jsonx.AppendArrayEnd(doc)
	doc = jsonx.AppendDocumentElementStart(doc, "props")
	for key, prop := range d.Props {
		doc = jsonx.AppendStringArrayElement(doc, key, prop)
	}
	doc = jsonx.AppendDocumentEnd(doc)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *ObjectInfo) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		k := cur.Key()
		switch k {
		case "id":
			d.Id = cur.Int32()
		case "parent":
			d.Parent = jsonx.DecodeInt32Array(cur.Value())
		case "class":
			d.Class = cur.String()
		case "shape":
			d.Shape = cur.String()
		case "coord":
			coord := cur.Value()
			d.Coord = make([][2]int32, coord.Size())
			for i := 0; coord.Next(); i++ {
				point := coord.Value()
				point.Next()
				d.Coord[i][0] = int32(point.Float())
				point.Next()
				d.Coord[i][1] = int32(point.Float())
			}
		case "props":
			sub := cur.Value()
			d.Props = make(map[string][]string, sub.Size())
			for sub.Next() {
				d.Props[sub.Key()] = jsonx.DecodeStringArray(sub.Value())
			}
		}
	}
}

func (d LabelTags) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		sub := cur.Value()
		temp := make(map[string][]string, sub.Size())
		for sub.Next() {
			temp[sub.Key()] = jsonx.DecodeStringArray(sub.Value())
		}
		d[cur.Key()] = temp
	}
}

func (d *LabelInfo) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "ver", d.Ver)
	doc = jsonx.AppendStringElement(doc, "fileName", d.FileName)
	doc = jsonx.AppendInt32Element(doc, "schemaVer", d.SchemaVer)
	doc = jsonx.AppendInt64Element(doc, "baseId", d.BaseId)
	doc = jsonx.AppendInt32Element(doc, "frameNo", d.FrameNum)
	doc = jsonx.AppendInt32ArrayElement(doc, "imgSize", []int32{d.ImgSize[0], d.ImgSize[1]})
	doc = jsonx.AppendDocumentElementStart(doc, "tags")
	for key, prop := range d.Tags {
		doc = jsonx.AppendDocumentElementStart(doc, key)
		for key2, prop2 := range prop {
			doc = jsonx.AppendStringArrayElement(doc, key2, prop2)
		}
		doc = jsonx.AppendDocumentEnd(doc)
	}
	doc = jsonx.AppendDocumentEnd(doc)
	doc = jsonx.AppendArrayElementStart(doc, "objects")
	for _, object := range d.Objects {
		doc = object.MarshalJson(doc)
	}
	doc = jsonx.AppendArrayEnd(doc)
	jsonx.AppendInt32Element(doc, "nextId", d.nextId)
	jsonx.AppendTimestampElement(doc, "updateTime", d.updateTime)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *LabelInfo) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		k := cur.Key()
		switch k {
		case "ver":
			d.Ver = cur.String()
		case "baseId":
			d.BaseId = cur.Int64()
		case "frameNo":
			d.FrameNum = cur.Int32()
		case "imgSize", "imageSize":
			imageSize := cur.Value()
			for i := 0; i < 2; i++ {
				imageSize.Next()
				d.ImgSize[i] = imageSize.Int32()
			}
		case "tags":
			sub := cur.Value()
			d.Tags = make(LabelTags, sub.Size())
			LabelTags(d.Tags).Decode(sub)
		case "objects":
			sub := cur.Value()
			d.Objects = make([]ObjectInfo, sub.Size())
			for i := 0; sub.Next(); i++ {
				d.Objects[i].Decode(sub.Value())
			}
		}
	}
}
