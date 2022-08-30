package jsonx

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

const (
	endField  = -1
	endArray  = -2
	endObject = -3
	endError  = -9
)

type marker struct {
	doc  []byte
	dp   int
	list []value
}

func (m *marker) readDoc() error {
	m.list = append(m.list, value{next: -1})
	return m.readValue()
}

func (m *marker) readValue() error {
	for ; m.dp < len(m.doc); m.dp++ {
		if m.doc[m.dp] > 0x20 {
			break
		}
	}
	switch m.doc[m.dp] {
	case '{':
		if err := m.readObject(); err != nil {
			return err
		}
	case '[':
		if err := m.readArray(); err != nil {
			return err
		}
	case 'n':
		if err := m.readNull(); err != nil {
			return err
		}
	case 't':
		if err := m.readTrue(); err != nil {
			return err
		}
	case 'f':
		if err := m.readFalse(); err != nil {
			return err
		}
	case '"':
		if err := m.readString(); err != nil {
			return err
		}
	default:
		return m.readNumber()
	}
	return nil
}

func (m *marker) readNull() error {
	v := &m.list[len(m.list)-1]
	v.kind = TypeNull
	v.valS = m.dp
	if m.doc[m.dp+1] != 'u' || m.doc[m.dp+2] != 'l' || m.doc[m.dp+3] != 'l' {
		return errors.New("error json value")
	}
	m.dp += 4
	v.valE = m.dp
	return nil
}

func (m *marker) readTrue() error {
	v := &m.list[len(m.list)-1]
	v.kind = TypeBool
	v.valS = m.dp
	if m.doc[m.dp+1] != 'r' || m.doc[m.dp+2] != 'u' || m.doc[m.dp+3] != 'e' {
		return errors.New("error json value")
	}
	m.dp += 4
	v.valE = m.dp
	return nil
}

func (m *marker) readFalse() error {
	v := &m.list[len(m.list)-1]
	v.kind = TypeBool
	v.valS = m.dp
	if m.doc[m.dp+1] != 'a' || m.doc[m.dp+2] != 'l' || m.doc[m.dp+3] != 's' || m.doc[m.dp+4] != 'e' {
		return errors.New("error json value")
	}
	m.dp += 5
	v.valE = m.dp
	return nil
}

func (m *marker) readNumber() error {
	v := &m.list[len(m.list)-1]
	v.kind = TypeInteger
	v.valS = m.dp
	pointCount := 0
	if m.doc[m.dp] == '-' {
		m.dp += 1
	}
	err := true
	for ; m.dp < len(m.doc); m.dp++ {
		if '0' <= m.doc[m.dp] && m.doc[m.dp] <= '9' {
			err = false
			continue
		}
		if m.doc[m.dp] <= 0x20 || m.doc[m.dp] == ',' || m.doc[m.dp] == ']' || m.doc[m.dp] == '}' {
			break
		}
		if m.doc[m.dp] == '.' && pointCount == 0 {
			v.kind = TypeFloat
			pointCount += 1
		} else {
			return errors.New("error json value")
		}
	}
	v.valE = m.dp
	if err {
		return errors.New("error json value")
	}
	return nil
}

func (m *marker) readString() error {
	m.dp += 1
	v := &m.list[len(m.list)-1]
	v.kind = TypeString
	v.valS = m.dp
	n, err := m.decodeString()
	if err != nil {
		return err
	}
	v.valE = v.valS + n
	return nil
}

func (m *marker) readArray() error {
	vp := len(m.list) - 1
	m.list[vp].kind = TypeArray
	m.list[vp].valS = m.dp
	m.dp += 1
	count := 0
	ch := m.readEndChar()
	if ch == endError {
		for {
			p := len(m.list)
			m.list = append(m.list, value{next: -1})
			if err := m.readValue(); err != nil {
				return err
			}
			count += 1
			ec := m.readEndChar()
			if ec == endField {
				m.list[p].next = len(m.list)
			} else if ec == endArray {
				break
			} else {
				return errors.New("array not closed")
			}
		}
	} else if ch != endArray {
		return errors.New("array not closed")
	}
	m.list[vp].valE = m.dp
	m.list[vp].len = count
	return nil
}

func (m *marker) readObject() error {
	vp := len(m.list) - 1
	m.list[vp].kind = TypeObject
	m.list[vp].valS = m.dp
	m.dp += 1
	count := 0
	ch := m.readEndChar()
	if ch == endError {
		for {
			p := len(m.list)
			sub := value{next: -1}
			if err := m.readKey(&sub); err != nil {
				return err
			}
			m.list = append(m.list, sub)
			if err := m.readValue(); err != nil {
				return err
			}
			count += 1
			ec := m.readEndChar()
			if ec == endField {
				m.list[p].next = len(m.list)
			} else if ec == endObject {
				break
			} else {
				return errors.New("object not closed")
			}
		}
	} else if ch != endObject {
		return errors.New("object not closed")
	}
	m.list[vp].valE = m.dp
	m.list[vp].len = count
	return nil
}

func (m *marker) readKey(v *value) error {
	for ; m.dp < len(m.doc); m.dp++ {
		if m.doc[m.dp] > 0x20 {
			if m.doc[m.dp] == '"' {
				m.dp += 1
				v.keyS = m.dp
				break
			}
			return errors.New("not quotation mark before key")
		}
	}
	if v.keyS == 0 {
		return errors.New("not find quotation mark")
	}
	n, err := m.decodeString()
	if err != nil {
		return err
	}
	v.keyE = v.keyS + n
	for ; m.dp < len(m.doc); m.dp++ {
		if m.doc[m.dp] > 0x20 {
			if m.doc[m.dp] == ':' {
				m.dp += 1
				return nil
			}
			break
		}
	}
	return errors.New("not colon after key")
}

func (m *marker) readEndChar() int {
	for ; m.dp < len(m.doc); m.dp++ {
		if m.doc[m.dp] > 0x20 {
			break
		}
	}
	switch m.doc[m.dp] {
	case ',':
		m.dp += 1
		return endField
	case ']':
		m.dp += 1
		return endArray
	case '}':
		m.dp += 1
		return endObject
	default:
		return endError
	}
}

func (m *marker) decodeString() (int, error) {
	newStr := m.doc[m.dp:m.dp]
	for ; m.dp < len(m.doc); m.dp++ {
		if m.doc[m.dp] == '"' {
			m.dp += 1
			return len(newStr), nil
		}
		if m.doc[m.dp] == '\\' {
			m.dp += 1
			if m.doc[m.dp] == 'u' {
				i, err := strconv.ParseUint(string(m.doc[m.dp+1:m.dp+5]), 16, 32)
				if err != nil {
					return 0, err
				}
				newStr = newStr[:len(newStr)+4]
				n := utf8.EncodeRune(newStr[len(newStr)-4:], rune(i))
				newStr = newStr[:len(newStr)+n-4]
				m.dp += 4
				continue
			}
		}
		newStr = append(newStr, m.doc[m.dp])
	}
	return 0, errors.New("not quotation mark after string value")
}
