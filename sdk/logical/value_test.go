package logical

import (
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewValue(t *testing.T) {

	t.Parallel()

	thouShallPass := func(in interface{}, kind byte, tst *testing.T) {
		val, err := NewValue(in, kind)
		assert.Nil(tst, err)
		assert.NotNil(tst, val)
		assert.Equal(tst, val.Kind, kind)
		assert.NotZero(tst, val.Size)

		switch kind {
		case Boolean:
			var test bool
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case String:
			var test string
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Int:
			var test int
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Int8:
			var test int8
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Int16:
			var test int16
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Int32:
			var test int32
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Int64:
			var test int64
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Uint:
			var test uint
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Uint8:
			var test uint8
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Uint16:
			var test uint16
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Uint32:
			var test uint32
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Uint64:
			var test uint64
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Float32:
			var test float32
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Float64:
			var test float64
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case BigFloat:
			var test *big.Float
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case BigInt:
			var test *big.Int
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Timestamp:
			var test time.Time
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		case Duration:
			var test time.Duration
			err = val.DecodeContent(&test)
			assert.Nil(tst, err)
			assert.NotNil(tst, test)
			assert.Equal(tst, in, test, fmt.Sprintf("expected = (%v), actual = (%v)", in, test))
		default:
			assert.Fail(tst, "something is wrong with the input kind")
		}

		encodedValue := val.Encode()
		assert.NotNil(tst, encodedValue)

		decodedValue, err := DecodeValue(encodedValue)
		assert.Nil(tst, err)
		assert.NotNil(tst, decodedValue)
		assert.Equal(tst, val, decodedValue)
		assert.EqualValues(tst, val, decodedValue)
	}

	thyKindShallNotPass := func(in interface{}, kind byte, tst *testing.T) {
		_, err := NewValue(in, kind)
		assert.NotNil(tst, err)
		assert.Contains(tst, err.Error(), "input type does not match kind")
	}

	t.Run("Boolean", func(t *testing.T) {
		t.Parallel()
		thouShallPass(true, Boolean, t)
		thouShallPass(false, Boolean, t)
		thyKindShallNotPass(true, Int, t)
		thyKindShallNotPass(false, BigInt, t)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		thouShallPass("hello world", String, t)
		thouShallPass("3.2", String, t)
		thouShallPass("10", String, t)
		thyKindShallNotPass("hello world", Boolean, t)
		thyKindShallNotPass("3.2", Float32, t)
		thyKindShallNotPass("10", Int8, t)
	})

	t.Run("Integers", func(t *testing.T) {
		t.Parallel()
		numbers := []int{0, 1, -1, 2, -2, 314, 159}
		t.Run("Int", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(number, Int, subtest)
				thyKindShallNotPass(number, String, subtest)
			}
		})
		t.Run("8 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(int8(number), Int8, subtest)
				thyKindShallNotPass(int8(number), Int, subtest)
			}
		})
		t.Run("16 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(int16(number), Int16, subtest)
				thyKindShallNotPass(int16(number), Int32, subtest)
			}
		})
		t.Run("32 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(int32(number), Int32, subtest)
				thyKindShallNotPass(int32(number), Int8, subtest)
			}
		})
		t.Run("64 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(int64(number), Int64, subtest)
				thyKindShallNotPass(int64(number), Boolean, subtest)
			}
		})
	})

	t.Run("Unsigned Integers", func(t *testing.T) {
		t.Parallel()
		numbers := []uint{0, 1, 314, 159, math.MaxUint8}
		t.Run("Uint", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(number, Uint, subtest)
				thyKindShallNotPass(number, String, subtest)
			}
		})
		t.Run("8 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(uint8(number), Uint8, subtest)
				thyKindShallNotPass(uint8(number), Uint, subtest)
			}
		})
		t.Run("16 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(uint16(number), Uint16, subtest)
				thyKindShallNotPass(uint16(number), Uint, subtest)
			}
		})
		t.Run("32 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(uint32(number), Uint32, subtest)
				thyKindShallNotPass(uint32(number), Uint, subtest)
			}
		})
		t.Run("64 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(uint64(number), Uint64, subtest)
				thyKindShallNotPass(uint64(number), Uint, subtest)
			}
		})
	})

	t.Run("Floating Point", func(t *testing.T) {
		t.Parallel()
		numbers := []float64{-12.3, 12.3, 15.1, 3.1415, 1e3, -0.15e7, 2e7}
		t.Run("32 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(float32(number), Float32, subtest)
			}
		})
		t.Run("64 bit", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range numbers {
				thouShallPass(number, Float64, subtest)
			}
		})
	})

	t.Run("Big Numbers", func(t *testing.T) {
		t.Parallel()
		ints := []int64{1, 2, 3, 100, 500, 2000}
		floats := []float64{-12.3, 12.3, 15.1, 3.1415, 1e3, -0.15e7, 2e7}
		t.Run("Integers", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range ints {
				thouShallPass(big.NewInt(number), BigInt, subtest)
			}
		})
		t.Run("Floats", func(subtest *testing.T) {
			subtest.Parallel()
			for _, number := range floats {
				thouShallPass(big.NewFloat(number), BigFloat, subtest)
			}
		})
	})

	t.Run("Timestamps", func(t *testing.T) {
		t.Parallel()
		timestamps := []time.Time{
			time.Now().UTC(),
			time.Date(2020, 10, 20, 1, 2, 3, 0, time.UTC),
			time.Unix(1619872786, 0),
			time.Unix(1619872786, 21312),
		}
		for _, elem := range timestamps {
			thouShallPass(elem, Timestamp, t)
			thyKindShallNotPass(elem, Int, t)
		}
	})

	t.Run("Durations", func(t *testing.T) {
		t.Parallel()
		durations := []time.Duration{
			3 * time.Second,
			1 * time.Nanosecond,
			100 * time.Hour,
			0 * time.Minute,
		}
		for _, elem := range durations {
			thouShallPass(elem, Duration, t)
			thyKindShallNotPass(elem, Int, t)
		}
	})

}

func Test_NegativeEncodeDecode(t *testing.T) {

	var value Value

	// Encoding nil content should fail
	err := value.EncodeContent(nil)
	assert.NotNil(t, err)

	// Creating a new value with nil content should fail
	// since type assertions will kick in
	_, err = NewValue(nil, Boolean)
	assert.NotNil(t, err)

	// Decoding into nil should also fail
	err = value.DecodeContent(nil)
	assert.NotNil(t, err)

	// Decoding nil input into a value should fail
	_, err = DecodeValue(nil)
	assert.NotNil(t, err)

}
