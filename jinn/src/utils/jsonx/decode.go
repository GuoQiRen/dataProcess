package jsonx

import (
	"errors"
	"strconv"
	"unsafe"
)

const (
	TypeNull = iota + 1
	TypeBool
	TypeInteger
	TypeFloat
	TypeString
	TypeArray
	TypeObject
)

var EmptyCursor = func() *Cursor {
	cur := Cursor{beg: value{next: -1}}
	cur.current = &cur.beg
	return &cur
}()

type Decoder interface {
	Decode(cur *Cursor)
}

type value struct {
	kind int
	keyS int
	keyE int
	valS int
	valE int
	len  int // 对象，数组长度
	next int
}

type Cursor struct {
	data    []byte  // json字符串
	list    []value // json对象
	beg     value   // 虚拟开始对象
	idx     int
	current *value
	size    int
}

func (c *Cursor) Clone() *Cursor {
	cur := *c
	cur.data = make([]byte, len(c.data))
	copy(cur.data, c.data)
	return &cur
}

func (c *Cursor) Reset() {
	c.current = &c.beg
}

func (c *Cursor) Next() bool {
	if c.current.next >= 0 {
		c.idx = c.current.next
		c.current = &c.list[c.idx]
		return true
	}
	return false
}

func (c *Cursor) Key() string {
	return string(c.data[c.current.keyS:c.current.keyE])
}

func (c *Cursor) Type() int {
	return c.current.kind
}

func (c *Cursor) Size() int {
	return c.size
}

func (c *Cursor) Value() *Cursor {
	cur := &Cursor{data: c.data, list: c.list, beg: *c.current, size: c.current.len}
	if c.current.len > 0 {
		cur.beg.next = c.idx + 1
	} else {
		cur.beg.next = -1
	}
	cur.current = &cur.beg
	return cur
}

func (c *Cursor) Bool() bool {
	return c.data[c.current.valS] == 't'
}

func (c *Cursor) Int32() int32 {
	buf := c.data[c.current.valS:c.current.valE]
	i, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&buf)), 10, 32)
	return int32(i)
}

func (c *Cursor) Uint32() uint32 {
	buf := c.data[c.current.valS:c.current.valE]
	i, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&buf)), 10, 32)
	return uint32(i)
}

func (c *Cursor) Int64() int64 {
	buf := c.data[c.current.valS:c.current.valE]
	i, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&buf)), 10, 64)
	return i
}

func (c *Cursor) Float() float64 {
	buf := c.data[c.current.valS:c.current.valE]
	i, _ := strconv.ParseFloat(*(*string)(unsafe.Pointer(&buf)), 64)
	return i
}

func (c *Cursor) String() string {
	return string(c.data[c.current.valS:c.current.valE])
}

func Parse(data []byte) (cur *Cursor, err error) {
	if len(data) == 0 {
		return nil, errors.New("json document is empty")
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("invalid json document")
		}
	}()
	m := marker{doc: data, list: make([]value, 0, 32)}
	if err = m.readDoc(); err != nil {
		return nil, err
	}
	cur = &Cursor{data: data, list: m.list, size: 1}
	cur.current = &m.list[0]
	return
}

func ParseObject(data []byte) (*Cursor, error) {
	cur, err := Parse(data)
	if err != nil {
		return nil, err
	}
	if cur.Type() != TypeObject {
		return nil, errors.New("not json object type")
	}
	return cur.Value(), nil
}

func DecodeObject(cur *Cursor) map[string]interface{} {
	result := make(map[string]interface{}, cur.Size())
	for cur.Next() {
		switch cur.Type() {
		case TypeNull:
			result[cur.Key()] = nil
		case TypeBool:
			result[cur.Key()] = cur.Bool()
		case TypeInteger:
			result[cur.Key()] = cur.Int64()
		case TypeFloat:
			result[cur.Key()] = cur.Float()
		case TypeString:
			result[cur.Key()] = cur.String()
		case TypeArray:
			result[cur.Key()] = DecodeArray(cur.Value())
		case TypeObject:
			result[cur.Key()] = DecodeObject(cur.Value())
		}
	}
	return result
}

func DecodeArray(cur *Cursor) []interface{} {
	result := make([]interface{}, 0, cur.Size())
	for cur.Next() {
		switch cur.Type() {
		case TypeNull:
			result = append(result, nil)
		case TypeBool:
			result = append(result, cur.Bool())
		case TypeInteger:
			result = append(result, cur.Int64())
		case TypeFloat:
			result = append(result, cur.Float())
		case TypeString:
			result = append(result, cur.String())
		case TypeArray:
			result = append(result, DecodeArray(cur.Value()))
		case TypeObject:
			result = append(result, DecodeObject(cur.Value()))
		}
	}
	return result
}

func DecodeStringArray(cur *Cursor) []string {
	arr := make([]string, 0, cur.Size())
	for cur.Next() {
		arr = append(arr, cur.String())
	}
	return arr
}

func DecodeInt64Array(cur *Cursor) []int64 {
	arr := make([]int64, 0, cur.Size())
	for cur.Next() {
		arr = append(arr, cur.Int64())
	}
	return arr
}

func DecodeInt32Array(cur *Cursor) []int32 {
	arr := make([]int32, 0, cur.Size())
	for cur.Next() {
		arr = append(arr, cur.Int32())
	}
	return arr
}
