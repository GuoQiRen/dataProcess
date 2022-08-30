package jsonx

import (
	"strconv"
	"time"
)

type Doc []byte

func (d Doc) String() string {
	return string(d[:len(d)-1])
}

func (d Doc) Bytes() []byte {
	return d[:len(d)-1]
}

type Marshaler interface {
	MarshalJson(doc Doc) Doc
}

func AppendHeader(dst Doc, key string) Doc {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"')
	return append(dst, ':')
}

func AppendDocumentStart(dst Doc) Doc {
	return append(dst, '{')
}

func AppendDocumentEnd(dst Doc) Doc {
	if dst[len(dst)-1] == ',' {
		dst[len(dst)-1] = '}'
	} else {
		dst = append(dst, '}')
	}
	return append(dst, ',')
}

func AppendArrayStart(dst Doc) Doc {
	return append(dst, '[')
}

func AppendArrayEnd(dst Doc) Doc {
	if dst[len(dst)-1] == ',' {
		dst[len(dst)-1] = ']'
	} else {
		dst = append(dst, ']')
	}
	return append(dst, ',')
}

func AppendDocumentElementStart(dst Doc, key string) Doc {
	return AppendDocumentStart(AppendHeader(dst, key))
}

func AppendArrayElementStart(dst Doc, key string) Doc {
	return AppendArrayStart(AppendHeader(dst, key))
}

func AppendString(dst Doc, s string) Doc {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			dst = append(dst, "\\\""...)
		case '\\':
			dst = append(dst, "\\\\"...)
		default:
			dst = append(dst, s[i])
		}
	}
	dst = append(dst, '"')
	return append(dst, ',')
}

func AppendStringElement(dst Doc, key string, s string) Doc {
	return AppendString(AppendHeader(dst, key), s)
}

func AppendInt64(dst Doc, i int64) Doc {
	dst = append(dst, '"')
	dst = append(dst, strconv.FormatInt(i, 10)...)
	dst = append(dst, '"')
	return append(dst, ',')
}

func AppendInt64Element(dst Doc, key string, i int64) Doc {
	return AppendInt64(AppendHeader(dst, key), i)
}

func AppendUint32(dst Doc, i uint32) Doc {
	dst = append(dst, strconv.FormatUint(uint64(i), 10)...)
	return append(dst, ',')
}

func AppendUint32Element(dst Doc, key string, i uint32) Doc {
	return AppendUint32(AppendHeader(dst, key), i)
}

func AppendInt32(dst Doc, i int32) Doc {
	dst = append(dst, strconv.FormatInt(int64(i), 10)...)
	return append(dst, ',')
}

func AppendInt32Element(dst Doc, key string, i int32) Doc {
	return AppendInt32(AppendHeader(dst, key), i)
}

func AppendBool(dst Doc, b bool) Doc {
	if b {
		return append(dst, "true,"...)
	} else {
		return append(dst, "false,"...)
	}
}

func AppendBoolElement(dst Doc, key string, b bool) Doc {
	return AppendBool(AppendHeader(dst, key), b)
}

func AppendNull(dst Doc) Doc {
	return append(dst, "null,"...)
}

func AppendNullElement(dst Doc, key string) Doc {
	return AppendNull(AppendHeader(dst, key))
}

func AppendTimestamp(dst Doc, t time.Time) Doc {
	dst = append(dst, strconv.FormatInt(t.Unix(), 10)...)
	return append(dst, ',')
}

func AppendTimestampElement(dst Doc, key string, t time.Time) Doc {
	return AppendTimestamp(AppendHeader(dst, key), t)
}

func AppendValue(dst Doc, v string) Doc {
	dst = append(dst, v...)
	return append(dst, ',')
}

func AppendValueElement(dst Doc, key string, v string) Doc {
	return AppendValue(AppendHeader(dst, key), v)
}

func AppendStringArrayElement(dst Doc, key string, arr []string) Doc {
	doc := AppendArrayElementStart(dst, key)
	for _, v := range arr {
		doc = AppendString(doc, v)
	}
	return AppendArrayEnd(doc)
}

func AppendInt64ArrayElement(dst Doc, key string, arr []int64) Doc {
	doc := AppendArrayElementStart(dst, key)
	for _, v := range arr {
		doc = AppendInt64(doc, v)
	}
	return AppendArrayEnd(doc)
}

func AppendInt32ArrayElement(dst Doc, key string, arr []int32) Doc {
	doc := AppendArrayElementStart(dst, key)
	for _, v := range arr {
		doc = AppendInt32(doc, v)
	}
	return AppendArrayEnd(doc)
}

func MergeDocumentElement(dst Doc, m Marshaler) Doc {
	end := len(dst) - 1
	bak := dst[end]
	dst = m.MarshalJson(dst[:end])
	dst[end] = bak
	dst[len(dst)-2] = ','
	return dst[:len(dst)-1]
}
