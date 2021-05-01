package logical

import (
	"fmt"
	"math/big"
	"time"

	"github.com/momentag/ebliss/sdk/helpers/gobutils"
)

const (
	Boolean = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	BigInt
	BigFloat
	String
	Timestamp
	Duration
)

type Value struct {
	Content   []byte
	Size      int
	Retention time.Duration
	Kind      byte
}

func AssertValue(in interface{}, kind byte) bool {
	switch in.(type) {
	case bool:
		return kind == Boolean
	case string:
		return kind == String
	case int:
		return kind == Int
	case int8:
		return kind == Int8
	case int16:
		return kind == Int16
	case int32:
		return kind == Int32
	case int64:
		return kind == Int64
	case uint:
		return kind == Uint
	case uint8:
		return kind == Uint8
	case uint16:
		return kind == Uint16
	case uint32:
		return kind == Uint32
	case uint64:
		return kind == Uint64
	case float32:
		return kind == Float32
	case float64:
		return kind == Float64
	case *big.Int:
		return kind == BigInt
	case *big.Float:
		return kind == BigFloat
	case time.Time:
		return kind == Timestamp
	case time.Duration:
		return kind == Duration
	default:
		return false
	}
}

func NewValue(in interface{}, kind byte) (Value, error) {
	var val Value
	if AssertValue(in, kind) {
		val.Kind = kind
		err := val.EncodeContent(in)
		if err != nil {
			return val, fmt.Errorf("cannot encode input (%v)", err.Error())
		} else {
			return val, nil
		}
	}
	return val, fmt.Errorf("input type does not match kind")
}

func DecodeValue(in []byte) (Value, error) {
	var value Value
	err := gobutils.DecodeData(in, &value)
	return value, err
}

func (v *Value) Encode() []byte {
	buf, _ := gobutils.EncodeInterface(v)
	return buf.Bytes()
}

func (v *Value) EncodeContent(in interface{}) error {
	if buf, err := gobutils.EncodeInterface(in); err != nil {
		return err
	} else {
		v.Size = buf.Len()
		v.Content = buf.Bytes()
		return nil
	}
}

func (v *Value) DecodeContent(out interface{}) error {
	return gobutils.DecodeData(v.Content, out)
}