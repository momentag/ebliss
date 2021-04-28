package parseutils

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseCapacityString(t *testing.T) {
	t.Parallel()

	// Test legit strings
	tests := map[string]uint64{
		"32.5kb":         32.5e3,
		"1.5MB":          1.5e6,
		"7GB":            7e9,
		"0.125kib":       128,
		"1":              1,
		"5b":             5,
		"0.000000005GB":  5,
		"0.000123456GiB": 132560,
		"-0.5kb":         0.5e3,
		"0.12mib":        125829,
		"8TB":            8e12,
		"0.0000015TiB":   1649267,
		"1.2e3kB":        1.2e6,
		"5.5e+3kB":       5.5e6,
	}

	for k, expected := range tests {
		actual, err := ParseCapacityString(k)
		assert.Nil(t, err, fmt.Sprintf("errored for k=%v and expected %v but got %v", k, expected, actual))
		assert.Equal(t, expected, actual, fmt.Sprintf("k=%v and expected %v but got %v", k, expected, actual))
	}

	testNumber := func(in interface{}, expected uint64) {
		actual, err := ParseCapacityString(in)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, fmt.Sprintf("in %v, expected %v, got %v", in, expected, actual))
	}

	// Test legit number cases
	testNumber(1, 1)
	testNumber(int8(3), 3)
	testNumber(int16(16), 16)
	testNumber(int32(32), 32)
	testNumber(int64(64), 64)
	testNumber(uint(1), 1)
	testNumber(uint8(15), 15)
	testNumber(uint16(16), 16)
	testNumber(uint32(32), 32)
	testNumber(uint64(64), 64)
	testNumber(float32(100.25), 100)
	testNumber(7.25, 7)

	// Test negative cases
	_, err := ParseCapacityString(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "input is nil")

	_, err = ParseCapacityString(true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not parse capacity from input")

	// Test input as JSON Number
	actual, err := ParseCapacityString(json.Number("3"))
	assert.Nil(t, err)
	assert.Equal(t, actual, uint64(3))

	// Test case where regex does not match
	_, err = ParseCapacityString("k3j8f")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not parse capacity from input")

	// Test empty string
	_, err = ParseCapacityString("")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "empty string provided, cannot parse capacity")

	// Test bananas range
	_, err = ParseCapacityString(fmt.Sprintf("%vTB", math.MaxFloat64))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "out of bounds")
}

func BenchmarkParseCapacityString_Strings(b *testing.B) {
	targets := []string{"32.5kB", "1.5MB", "7GB", "0.125kib", "1", "5b",
		"0.000000005GB", "-0.5kb", "0.12mib", "8TB", "0.0000015TiB", "1.2e3kB",
		"5.5e+3kB"}
	for target := range targets {
		for n := 0; n < b.N; n++ {
			_, _ = ParseCapacityString(target)
		}
	}
}

func TestParseDurationSecond(t *testing.T) {
	t.Parallel()

	// Test legit strings
	tests := map[string]time.Duration{
		"12s":   12 * time.Second,
		"1m":    60 * time.Second,
		"7m":    7 * time.Minute,
		"475s":  475 * time.Second,
		"1.5m":  1*time.Minute + 30*time.Second,
		"800ms": 800 * time.Millisecond,
	}

	for k, expected := range tests {
		actual, err := ParseDurationSecond(k)
		assert.Nil(t, err, fmt.Sprintf("errored for k=%v and expected %v but got %v", k, expected, actual))
		assert.Equal(t, expected, actual, fmt.Sprintf("k=%v and expected %v but got %v", k, expected, actual))
	}

	testNumber := func(in interface{}, expected time.Duration) {
		actual, err := ParseDurationSecond(in)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, fmt.Sprintf("in %v, expected %v, got %v", in, expected, actual))
	}

	// Test legit number cases
	testNumber(1, 1*time.Second)
	testNumber(int8(3), 3*time.Second)
	testNumber(int16(16), 16*time.Second)
	testNumber(int32(32), 32*time.Second)
	testNumber(int64(64), 64*time.Second)
	testNumber(uint(1), 1*time.Second)
	testNumber(uint8(15), 15*time.Second)
	testNumber(uint16(16), 16*time.Second)
	testNumber(uint32(32), 32*time.Second)
	testNumber(uint64(64), 64*time.Second)
	testNumber(float32(100.25), 100*time.Second)
	testNumber(7.25, 7*time.Second)

	// Test negative cases
	_, err := ParseDurationSecond(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "input is nil")

	_, err = ParseDurationSecond(true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "boolean value")

	// Test input as JSON Number
	actual, err := ParseDurationSecond(json.Number("3"))
	assert.Nil(t, err)
	assert.Equal(t, actual, 3*time.Second)

	// Test case where regex does not match
	_, err = ParseDurationSecond("k3j8f")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")

	// Test empty string
	_, err = ParseDurationSecond("")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "input string is empty")

	// Test time duration
	expected := 3 * time.Second
	actual, err = ParseDurationSecond(3 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}

func TestParseTime(t *testing.T) {
	t.Parallel()

	// Test legit strings
	tests := []struct {
		inp      interface{}
		valid    bool
		expected time.Time
	}{
		{
			"2020-12-11T09:08:07.654321Z",
			true,
			time.Date(2020, 12, 11, 9, 8, 7, 654321000, time.UTC),
		},
		{
			"2020-12-11T09:08:07+02:00",
			true,
			time.Date(2020, 12, 11, 7, 8, 7, 0, time.UTC),
		},
		{
			"2021-12-11T09:08:07Z",
			true,
			time.Date(2021, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			"2021-12-11T09:08:07",
			false,
			time.Time{},
		},
		{
			"1670749687",
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			1670749687,
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			uint32(1670749687),
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			json.Number("1670749687"),
			true,
			time.Date(2022, 12, 11, 9, 8, 7, 0, time.UTC),
		},
		{
			nil,
			true,
			time.Time{},
		},
		{
			struct{}{},
			false,
			time.Time{},
		},
		{
			true,
			false,
			time.Time{},
		},
	}

	for _, tc := range tests {
		actual, err := ParseTime(tc.inp)
		if tc.valid {
			assert.Nil(t, err, fmt.Sprintf("errored for k=%v and expected %v but got %v", tc.inp, tc.expected, actual))
			assert.Equal(t, tc.expected.UTC(), actual.UTC(), fmt.Sprintf("k=%v and expected %v but got %v", tc.inp, tc.expected, actual))
		} else {
			assert.NotNil(t, err)
		}
	}

}
