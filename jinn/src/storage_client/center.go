package client

import "dataProcess/jinn/src/utils/jsonx"

const (
	MethodCGetNode = "center.get_node"
	MethodNAuth    = "node.auth"
)

type NodeAuthInfo struct {
	Addr string
	Key  string
}

type CenterGetNodeParams struct {
	User     string
	Password string
	Token    string
	Class    uint32
}

func (p *CenterGetNodeParams) MarshalJson(doc jsonx.Doc) jsonx.Doc {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "user", p.User)
	doc = jsonx.AppendStringElement(doc, "password", p.Password)
	doc = jsonx.AppendStringElement(doc, "token", p.Token)
	doc = jsonx.AppendUint32Element(doc, "class", p.Class)
	return jsonx.AppendDocumentEnd(doc)
}

type CenterGetNodeResult struct {
	Info NodeAuthInfo
}

func (r *CenterGetNodeResult) Decode(cur *jsonx.Cursor) {
	for cur.Next() {
		switch cur.Key() {
		case "addr":
			r.Info.Addr = cur.String()
		case "key":
			r.Info.Key = cur.String()
		}
	}
}
