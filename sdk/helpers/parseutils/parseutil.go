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
	case int:
		return uint64(inp), nil
	case int8:
		return uint64(inp), nil
	case int16:
		return uint64(inp), nil
	case int32:
		return uint64(inp), nil
	case int64:
		return uint64(inp), nil
	case uint:
		return uint64(inp), nil
	case uint8:
		return uint64(inp), nil
	case uint16:
		return uint64(inp), nil
	case uint32:
		return uint64(inp), nil
	case uint64:
		return inp, nil
	case float32:
		return uint64(math.Round(float64(inp))), nil
	case float64:
		return uint64(math.Round(inp)), nil
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
	case string:
		return parseString(inp)
	case int:
		return time.Duration(inp) * time.Second, nil
	case int8:
		return time.Duration(inp) * time.Second, nil
	case int16:
		return time.Duration(inp) * time.Second, nil
	case int32:
		return time.Duration(inp) * time.Second, nil
	case int64:
		return time.Duration(inp) * time.Second, nil
	case uint:
		return time.Duration(inp) * time.Second, nil
	case uint8:
		return time.Duration(inp) * time.Second, nil
	case uint16:
		return time.Duration(inp) * time.Second, nil
	case uint32:
		return time.Duration(inp) * time.Second, nil
	case uint64:
		return time.Duration(inp) * time.Second, nil
	case float32:
		return time.Duration(inp) * time.Second, nil
	case float64:
		return time.Duration(inp) * time.Second, nil
	case time.Duration:
		return inp, nil
	default:
		return 0, fmt.Errorf("could not parse duration from input %v", in)
	}
}
