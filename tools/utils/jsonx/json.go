package jsonx

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

func Marshal(i interface{}, doc []byte) []byte {
	return marshalValue(doc, reflect.TypeOf(i), reflect.ValueOf(i)).Bytes()
}

func marshalValue(doc Doc, tp reflect.Type, value reflect.Value) Doc {
	if !value.IsValid() {
		return AppendNull(doc)
	}
	switch tp.Kind() {
	case reflect.Bool:
		return AppendBool(doc, value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return AppendInt32(doc, int32(value.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return AppendUint32(doc, uint32(value.Uint()))
	case reflect.Int64:
		return AppendInt64(doc, value.Int())
	case reflect.String:
		return AppendString(doc, value.String())
	case reflect.Array, reflect.Slice:
		doc = AppendArrayStart(doc)
		et := tp.Elem()
		for i := 0; i < value.Len(); i++ {
			doc = marshalValue(doc, et, value.Index(i))
		}
		return AppendArrayEnd(doc)
	case reflect.Map:
		doc = AppendDocumentStart(doc)
		iter := value.MapRange()
		et := tp.Elem()
		for iter.Next() {
			doc = AppendHeader(doc, iter.Key().String())
			doc = marshalValue(doc, et, iter.Value())
		}
		return AppendDocumentEnd(doc)
	case reflect.Struct:
		if tp.String() == "time.Time" {
			return AppendTimestamp(doc, value.Interface().(time.Time))
		}
		doc = AppendDocumentStart(doc)
		for i := 0; i < value.NumField(); i++ {
			elem := value.Field(i)
			if !elem.CanInterface() {
				continue
			}
			sf := tp.Field(i)
			name := sf.Tag.Get("json")
			if name == "-" {
				continue
			}
			if len(name) > 0 {
				doc = AppendHeader(doc, name)
			} else {
				doc = AppendHeader(doc, sf.Name)
			}
			doc = marshalValue(doc, sf.Type, elem)
		}
		return AppendDocumentEnd(doc)
	case reflect.Interface:
		if value.IsNil() {
			return AppendNull(doc)
		}
		elem := value.Elem()
		return marshalValue(doc, elem.Type(), elem)
	case reflect.Ptr:
		if value.IsNil() {
			return AppendNull(doc)
		}
		if value.CanInterface() && tp.Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()) {
			return value.Interface().(Marshaler).MarshalJson(doc)
		}
		return marshalValue(doc, tp.Elem(), value.Elem())
	default:
		return AppendNull(doc)
	}
}

func Unmarshal(doc []byte, i interface{}) error {
	tp := reflect.TypeOf(i)
	if tp.Kind() != reflect.Ptr {
		return errors.New("invalid argument")
	}
	cur, err := Parse(doc)
	if err != nil {
		return err
	}
	cur.Next()
	unmarshalValue(cur, tp, reflect.ValueOf(i))
	return nil
}

func unmarshalValue(cur *Cursor, tp reflect.Type, value reflect.Value) {
	for tp.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		if value.CanInterface() && tp.Implements(reflect.TypeOf((*Decoder)(nil)).Elem()) {
			value.Interface().(Decoder).Decode(cur.Value())
			return
		}
		tp = tp.Elem()
		value = value.Elem()
	}
	switch cur.Type() {
	case TypeBool:
		if tp.Kind() == reflect.Bool {
			value.SetBool(cur.Bool())
		}
	case TypeInteger:
		switch tp.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value.SetInt(cur.Int64())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			value.SetUint(uint64(cur.Uint32()))
		default:
			if tp.String() == "time.Time" {
				value.Set(reflect.ValueOf(time.Unix(cur.Int64(), 0)))
			}
		}
	case TypeFloat:
		switch tp.Kind() {
		case reflect.Float32, reflect.Float64:
			value.SetFloat(cur.Float())
		}
	case TypeString:
		if tp.Kind() == reflect.String {
			value.SetString(cur.String())
		}
	case TypeArray:
		unmarshalArray(cur.Value(), tp, value)
	case TypeObject:
		if value.CanAddr() {
			pv := value.Addr()
			if pv.CanInterface() && pv.Type().Implements(reflect.TypeOf((*Decoder)(nil)).Elem()) {
				pv.Interface().(Decoder).Decode(cur.Value())
				return
			}
		}
		unmarshalObject(cur.Value(), tp, value)
	}
}

func unmarshalArray(cur *Cursor, tp reflect.Type, value reflect.Value) {
	var pv *reflect.Value
	if tp.Kind() == reflect.Slice {
		temp := reflect.MakeSlice(tp, cur.Size(), cur.Size())
		pv = &temp
	} else if tp.Kind() == reflect.Array {
		if value.Len() != cur.Size() {
			return
		}
		pv = &value
	} else {
		return
	}
	et := tp.Elem()
	for i := 0; cur.Next(); i++ {
		unmarshalValue(cur, et, pv.Index(i))
	}
	if tp.Kind() == reflect.Slice {
		value.Set(*pv)
	}
}

func unmarshalObject(cur *Cursor, tp reflect.Type, value reflect.Value) {
	switch tp.Kind() {
	case reflect.Map:
		kt := tp.Key()
		et := tp.Elem()
		if kt.Kind() == reflect.String {
			if et.Kind() == reflect.Interface {
				data := DecodeObject(cur)
				value.Set(reflect.ValueOf(data))
				return
			}
			for cur.Next() {
				elem := reflect.New(et).Elem()
				unmarshalValue(cur, et, elem)
				value.SetMapIndex(reflect.ValueOf(cur.Key()), elem)
			}
		}
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			sf := tp.Field(i)
			name := sf.Tag.Get("json")
			if name == "-" {
				continue
			}
			elem := value.Field(i)
			if !elem.CanSet() {
				continue
			}
			for cur.Next() {
				if cur.Key() == name || strings.EqualFold(cur.Key(), sf.Name) {
					unmarshalValue(cur, sf.Type, elem)
					break
				}
			}
			cur.Reset()
		}
	}
}
