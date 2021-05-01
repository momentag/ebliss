package logical

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/momentag/ebliss/sdk/helpers/gobutils"
	"github.com/stretchr/testify/assert"
)

func TestParseValue(t *testing.T) {
	kind := byte(String)

	vbl := Variable{
		Name:       "variable",
		Implements: kind,
	}

	t.Run("matching kind value should be correctly parsed", func(t *testing.T) {
		value := Value{
			Content:   nil,
			Size:      0,
			Retention: 0,
			Kind:      kind,
		}
		entry, ok := parseValue(value, &vbl)
		assert.True(t, ok)
		assert.NotNil(t, entry)
	})

	t.Run("value with wrong kind should be rejected", func(t *testing.T) {
		value := Value{
			Content:   nil,
			Size:      0,
			Retention: 0,
			Kind:      Boolean,
		}
		entry, ok := parseValue(value, &vbl)
		assert.False(t, ok)
		assert.NotNil(t, entry)
	})
}

func TestParseValueContent(t *testing.T) {
	vbl := Variable{
		Name:       "hello",
		Implements: String,
	}
	data := "world"
	wrongData := 1

	t.Run("appropriate value should be encoded into an entry", func(t *testing.T) {
		entry, ok := parseContent(data, &vbl)
		assert.True(t, ok)
		assert.NotNil(t, entry)
		assert.Equal(t, entry.Key, vbl)
		assert.NotNil(t, entry.Value)
	})

	t.Run("wrong data type should not be encoded into an entry", func(t *testing.T) {
		entry, ok := parseContent(wrongData, &vbl)
		assert.False(t, ok)
		assert.NotNil(t, entry)
	})
}

func TestParseBytes(t *testing.T) {

	vbl := Variable{
		Name:       "hello",
		Implements: String,
	}

	buf, _ := gobutils.EncodeInterface("world")
	content := buf.Bytes()

	val := Value{
		Content:   content,
		Size:      buf.Len(),
		Retention: 0,
		Kind:      String,
	}

	t.Run("encoded Value should be parsed", func(t *testing.T) {
		valbytes := val.Encode()
		entry, ok := parseBytes(valbytes, &vbl)
		assert.True(t, ok)
		assert.NotNil(t, entry)
		assert.Equal(t, val, entry.Value)
	})

	t.Run("encoded data with appropriate type should be parsed", func(t *testing.T) {

		test := func(buffer *bytes.Buffer, kind int) {
			content := buffer.Bytes()
			vbl := Variable{
				Name:       "variable",
				Implements: byte(kind),
			}
			val := Value{
				Content:   content,
				Size:      buffer.Len(),
				Retention: 0,
				Kind:      byte(kind),
			}
			entry, ok := parseBytes(content, &vbl)
			assert.True(t, ok)
			assert.NotNil(t, entry)
			assert.Equal(t, val, entry.Value, fmt.Sprintf("Kind %v", kind))
		}

		t.Run("Types", func(t *testing.T) {
			t.Parallel()
			list := map[int]interface{}{
				Boolean:   true,
				String:    "hello world",
				Int:       1,
				Int8:      int8(1),
				Int16:     int16(1),
				Int32:     int32(1),
				Int64:     int64(1),
				Uint:      uint(1),
				Uint8:     uint8(1),
				Uint16:    uint16(1),
				Uint32:    uint32(1),
				Uint64:    uint64(1),
				Float32:   float32(1.0),
				Float64:   1.0,
				BigInt:    big.NewInt(1),
				BigFloat:  big.NewFloat(1.0),
				Timestamp: time.Now().UTC(),
				Duration:  3 * time.Second,
			}
			for k, v := range list {
				t.Run(fmt.Sprintf("%v", k), func(t *testing.T) {
					buf, _ := gobutils.EncodeInterface(v)
					test(buf, k)
				})
			}
		})
	})

}

func TestNewEntry(t *testing.T) {

	t.Run("individual values", func(t *testing.T) {
		t.Parallel()
		list := map[int]interface{}{
			Boolean:   true,
			String:    "hello world",
			Int:       1,
			Int8:      int8(1),
			Int16:     int16(1),
			Int32:     int32(1),
			Int64:     int64(1),
			Uint:      uint(1),
			Uint8:     uint8(1),
			Uint16:    uint16(1),
			Uint32:    uint32(1),
			Uint64:    uint64(1),
			Float32:   float32(1.0),
			Float64:   1.0,
			BigInt:    big.NewInt(1),
			BigFloat:  big.NewFloat(1.0),
			Timestamp: time.Now().UTC(),
			Duration:  3 * time.Second,
		}
		for k, v := range list {
			t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
				vbl := Variable{
					Name:       "variable",
					Implements: byte(k),
				}
				entry, err := vbl.NewEntry(v)
				assert.Nil(t, err)
				assert.NotNil(t, entry)
				assert.Equal(t, entry.Key, vbl)
			})
		}
	})

}
