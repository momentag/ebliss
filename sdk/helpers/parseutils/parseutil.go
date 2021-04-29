package parseutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validCapacityString = regexp.MustCompile("^[\t ]*([-0-9.e+]+)[\t ]?([kmgtKMGT]?[iI]?[bB])?[\t ]*$")

func ParseInt64(in interface{}) (int64, error) {
	var ret int64
	if in == nil {
		return ret, errors.New("input is nil")
	}
	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}
	switch in.(type) {
	case string:
		inp := in.(string)
		if inp == "" {
			return 0, nil
		}
		left, err := strconv.ParseInt(inp, 10, 64)
		if err != nil {
			return ret, err
		}
		ret = left
	case int:
		ret = int64(in.(int))
	case int8:
		ret = int64(in.(int8))
	case int16:
		ret = int64(in.(int16))
	case int32:
		ret = int64(in.(int32))
	case int64:
		ret = in.(int64)
	case uint:
		ret = int64(in.(uint))
	case uint8:
		ret = int64(in.(uint8))
	case uint16:
		ret = int64(in.(uint16))
	case uint32:
		ret = int64(in.(uint32))
	case uint64:
		ret = int64(in.(uint64))
	case float32:
		ret = int64(math.Round(float64(in.(float32))))
	case float64:
		ret = int64(math.Round(in.(float64)))
	default:
		return 0, errors.New("could not parse int from input")
	}
	return ret, nil
}

func ParseUint64(in interface{}) (uint64, error) {
	var ret uint64
	if in == nil {
		return ret, errors.New("input is nil")
	}
	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}
	switch in.(type) {
	case string:
		inp := in.(string)
		if inp == "" {
			return 0, nil
		}
		left, err := strconv.ParseUint(inp, 10, 64)
		if err != nil {
			return ret, err
		}
		ret = left
	case int:
		ret = uint64(in.(int))
	case int8:
		ret = uint64(in.(int8))
	case int16:
		ret = uint64(in.(int16))
	case int32:
		ret = uint64(in.(int32))
	case int64:
		ret = uint64(in.(int64))
	case uint:
		ret = uint64(in.(uint))
	case uint8:
		ret = uint64(in.(uint8))
	case uint16:
		ret = uint64(in.(uint16))
	case uint32:
		ret = uint64(in.(uint32))
	case uint64:
		ret = in.(uint64)
	case float32:
		ret = uint64(math.Round(float64(in.(float32))))
	case float64:
		ret = uint64(math.Round(in.(float64)))
	default:
		return 0, errors.New("could not parse int from input")
	}
	return ret, nil
}

func ParseCapacityString(in interface{}) (uint64, error) {
	var result uint64

	if in == nil {
		return result, fmt.Errorf("input is nil")
	}

	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}

	parseString := func(input string) (uint64, error) {
		if input == "" {
			return result, fmt.Errorf("empty string provided, cannot parse capacity")
		}
		matches := validCapacityString.FindStringSubmatch(input)
		if len(matches) <= 1 {
			return result, fmt.Errorf("could not parse capacity from input - %v", matches)
		}
		var multiplier float64 = 1
		switch strings.ToLower(matches[2]) {
		case "kb":
			multiplier = 1e3
		case "kib":
			multiplier = 1024
		case "mb":
			multiplier = 1e6
		case "mib":
			multiplier = math.Pow(1024, 2)
		case "gb":
			multiplier = 1e9
		case "gib":
			multiplier = math.Pow(1024, 3)
		case "tb":
			multiplier = 1e12
		case "tib":
			multiplier = math.Pow(1024, 4)
		}
		size, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return result, err
		}
		if math.Abs(size) < math.MaxFloat64 {
			return uint64(math.Round(math.Abs(size) * multiplier)), nil
		} else {
			return uint64(0), fmt.Errorf("resulting number out of bounds")
		}
	}

	switch inp := in.(type) {
	case string:
		return parseString(inp)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return ParseUint64(inp)
	default:
		return result, errors.New("could not parse capacity from input")
	}
}

func ParseDurationSecond(in interface{}) (time.Duration, error) {
	var dur time.Duration

	if in == nil {
		return dur, fmt.Errorf("input is nil")
	}

	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}

	parseString := func(input string) (time.Duration, error) {
		if input == "" {
			return dur, fmt.Errorf("input string is empty")
		}
		if strings.HasSuffix(input, "s") || strings.HasSuffix(input, "m") ||
			strings.HasSuffix(input, "h") || strings.HasSuffix(input, "ms") {
			return time.ParseDuration(input)
		} else {
			secs, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return dur, err
			}
			return time.Duration(secs) * time.Second, nil
		}
	}

	switch inp := in.(type) {
	case bool:
		return 0, fmt.Errorf("cannot parse duration from boolean value")
	case string:
		return parseString(inp)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		val, err := ParseInt64(inp)
		if err != nil {
			return dur, err
		}
		return time.Duration(val) * time.Second, nil
	case time.Duration:
		return inp, nil
	default:
		return 0, fmt.Errorf("could not parse duration from input %v", in)
	}
}

func ParseTime(in interface{}) (time.Time, error) {
	var t time.Time

	if in == nil {
		return time.Time{}, nil
	}

	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}

	parseString := func(input string) (time.Time, error) {
		if input == "" {
			return t, fmt.Errorf("input string is empty")
		}
		var err error
		t, err := time.Parse(time.RFC3339Nano, input)
		if err == nil {
			return t, nil
		}
		t, err = time.Parse(time.RFC3339, input)
		if err == nil {
			return t, nil
		}
		epoch, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			return time.Unix(epoch, 0), nil
		}
		return t, errors.New("could not parse string as datetime")
	}

	switch inp := in.(type) {
	case string:
		return parseString(inp)
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64:
		intValue, err := ParseInt64(inp)
		if err != nil {
			return t, err
		}
		return time.Unix(intValue, 0), nil
	default:
		return t, errors.New("could not parse datetime")
	}
}
