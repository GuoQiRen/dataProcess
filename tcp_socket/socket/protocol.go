package socket

import (
	"dataProcess/tools/utils/jsonx"
)

const ErrorUnknownMethod = 9998
const ErrorUnknown = 9999

type JsonPacket struct {
	Id     uint32
	Method string
	Params *jsonx.Cursor
	Result *jsonx.Cursor
	Error  struct {
		Code    int32
		Message string
	}
	Object *jsonx.Cursor
}

func (p *JsonPacket) Unmarshal(doc []byte) error {
	cur, err := jsonx.ParseObject(doc)
	if err != nil {
		return err
	}
	for cur.Next() {
		switch cur.Key() {
		case "id":
			p.Id = cur.Uint32()
		case "method":
			p.Method = cur.String()
		case "params":
			p.Params = cur.Value()
		case "Result":
			p.Result = cur.Value()
		case "error":
			e := cur.Value()
			for e.Next() {
				switch e.Key() {
				case "code":
					p.Error.Code = e.Int32()
				case "message":
					p.Error.Message = e.String()
				}
			}
		}
	}
	p.Object = cur
	return nil
}

func MarshalError(id uint32, code int32, msg string, doc []byte) []byte {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "id", id)
	doc = jsonx.AppendDocumentElementStart(doc, "error")
	doc = jsonx.AppendInt32Element(doc, "code", code)
	doc = jsonx.AppendStringElement(doc, "message", msg)
	doc = jsonx.AppendDocumentEnd(doc)
	return jsonx.AppendDocumentEnd(doc)
}

func MarshalRespond(id uint32, reply jsonx.Marshaler, doc []byte) []byte {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendUint32Element(doc, "id", id)
	doc = jsonx.AppendHeader(doc, "Result")
	doc = reply.MarshalJson(doc)
	return jsonx.AppendDocumentEnd(doc)
}

type registerArgs struct {
	id string
}

func (a *registerArgs) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "id", a.id)
	return jsonx.AppendDocumentEnd(doc)
}

func (a *registerArgs) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "id":
			a.id = cur.String()
		}
	}
}

type EmptyReply struct {
}

func (r *EmptyReply) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	return jsonx.AppendNull(doc)
}

func (r *EmptyReply) Decode(_ *jsonx.Cursor) {

}

type DownloadReply struct {
	Result   bool
	ErrMsg   string
	JinnPath string
}

func (d *DownloadReply) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendBoolElement(doc, "Result", d.Result)
	doc = jsonx.AppendStringElement(doc, "errMsg", d.ErrMsg)
	doc = jsonx.AppendStringElement(doc, "jinnPath", d.JinnPath)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *DownloadReply) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "Result":
			d.Result = cur.Bool()
		case "errMsg":
			d.ErrMsg = cur.String()
		case "jinnPath":
			d.JinnPath = cur.String()
		}
	}
}

type TaskReqs struct {
	Id     string
	UserId string
}

func (d *TaskReqs) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "id", d.Id)
	doc = jsonx.AppendStringElement(doc, "userId", d.UserId)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *TaskReqs) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "id":
			d.Id = cur.String()
		case "userId":
			d.UserId = cur.String()
		}
	}
}

type TaskReply struct {
	Result  bool
	ErrMsg  string
	TaskId  string
	EndTime string
}

func (d *TaskReply) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendBoolElement(doc, "result", d.Result)
	doc = jsonx.AppendStringElement(doc, "taskId", d.TaskId)
	doc = jsonx.AppendStringElement(doc, "errMsg", d.ErrMsg)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *TaskReply) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "result":
			d.Result = cur.Bool()
		case "taskId":
			d.TaskId = cur.String()
		case "errMsg":
			d.ErrMsg = cur.String()
		}
	}
}

type CancelArgs struct {
	Cause string
}

func (a *CancelArgs) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "Cause", a.Cause)
	return jsonx.AppendDocumentEnd(doc)
}

type CancelReply struct {
	Result bool
}

func (r *CancelReply) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	return jsonx.AppendBool(doc, true)
}

func (r *CancelReply) Decode(cur *jsonx.Cursor) {
	r.Result = cur.Bool()
}

type ProcessReply struct {
	TaskId    string
	Status    string
	Stage     string
	Exception string
	EndTime   string
}

func (d *ProcessReply) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "taskId", d.TaskId)
	doc = jsonx.AppendStringElement(doc, "status", d.Status)
	doc = jsonx.AppendStringElement(doc, "stage", d.Stage)
	doc = jsonx.AppendStringElement(doc, "exception", d.Exception)
	doc = jsonx.AppendStringElement(doc, "endTime", d.EndTime)
	return jsonx.AppendDocumentEnd(doc)
}

func (d *ProcessReply) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "taskId":
			d.TaskId = cur.String()
		case "status":
			d.Status = cur.String()
		case "stage":
			d.Stage = cur.String()
		case "exception":
			d.Exception = cur.String()
		case "endTime":
			d.EndTime = cur.String()
		}
	}
}
