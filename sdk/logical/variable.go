package logical

import (
	"fmt"
	"math/big"
	"time"

	"github.com/momentag/ebliss/sdk/helpers/gobutils"
)

type Variable struct {
	Name       string
	Implements byte
}

// parseValue parses the Value in, giving a new Entry
// while at the same time guaranteeing type uniformity
func parseValue(in Value, v *Variable) (Entry, bool) {
	if v.Implements == in.Kind {
		return Entry{
			Key:   *v,
			Value: in,
		}, true
	} else {
		return Entry{}, false
	}
}

// parseContent parses the contents of an interface
// into a new Value, then combine that with the Variable to create a new Entry
func parseContent(in interface{}, v *Variable) (Entry, bool) {
	if val, err := NewValue(in, v.Implements); err == nil {
		return parseValue(val, v)
	} else {
		return Entry{}, false
	}
}

// parseBytes this is a special case, where we have two possibilities
// either we feed an encoded Value which unmarshalls into a Value struct
// or we are parsing Value contents
func parseBytes(in []byte, v *Variable) (Entry, bool) {
	if value, err := DecodeValue(in); err == nil {
		return parseValue(value, v)
	} else {
		var err error
		var entry Entry
		var ok bool
		switch v.Implements {
		case Boolean:
			var content bool
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case String:
			var content string
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Int:
			var content int
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Int8:
			var content int8
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Int16:
			var content int16
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Int32:
			var content int32
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Int64:
			var content int64
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Uint:
			var content uint
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Uint8:
			var content uint8
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Uint16:
			var content uint16
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Uint32:
			var content uint32
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Uint64:
			var content uint64
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Float32:
			var content float32
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Float64:
			var content float64
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case BigInt:
			var content *big.Int
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case BigFloat:
			var content *big.Float
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Timestamp:
			var content time.Time
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		case Duration:
			var content time.Duration
			err = gobutils.DecodeData(in, &content)
			entry, ok = parseContent(content, v)
		}
		if err == nil {
			return entry, ok
		} else {
			return Entry{}, false
		}
	}
}

func (v *Variable) NewEntry(in interface{}) (Entry, error) {
	if in == nil {
		return Entry{}, fmt.Errorf("input is nil")
	}
	switch input := in.(type) {
	case Value:
		if entry, ok := parseValue(input, v); ok {
			return entry, nil
		}
	case []byte:
		if entry, ok := parseBytes(input, v); ok {
			return entry, nil
		}
	default:
		if entry, ok := parseContent(in, v); ok {
			return entry, nil
		}
	}
	return Entry{}, fmt.Errorf("could not parse input to an entry")
}
