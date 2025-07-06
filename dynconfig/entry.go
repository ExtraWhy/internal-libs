package dynconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Entry struct {
	err      error
	path     string
	entry    interface{}
	entryMap map[string]interface{}
}

func NewEntry(entry interface{}) *Entry {
	return NewEntryWithRoot(entry, "/", nil)
}

func NewEntryWithRoot(entry interface{}, rootName string, err error) *Entry {
	return &Entry{
		err:      err,
		path:     fmt.Sprintf("[%s]", rootName),
		entry:    entry,
		entryMap: nil,
	}
}

func (en *Entry) Get(key string) *Entry {
	path := fmt.Sprintf("%s/%s", en.path, key)
	if en.err != nil {
		return &Entry{en.err, path, nil, nil}
	}

	if entryMap, err := en.AsMap(); err == nil {
		if nestedEntry, ok := entryMap[key]; ok {
			return &Entry{nil, path, nestedEntry, nil}
		}
		return &Entry{errors.New("element not found"), path, nil, nil}
	}
	return &Entry{errors.New("element not found/reachable (not a JSON/Map)"), path, nil, nil}
}

func (en *Entry) Exists() bool {
	return en.entry != nil
}

func (en *Entry) AsString() (string, error) {
	if en.err != nil {
		return "", fmt.Errorf("can't get element (%s) AsString : %w", en.path, en.err)
	}

	switch value := en.entry.(type) {
	case float32:
		return strconv.FormatFloat(float64(value), 'g', -1, 64), nil
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64), nil
	case int:
		return strconv.FormatInt(int64(value), 10), nil
	case int32:
		return strconv.FormatInt(int64(value), 10), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		return value, nil
	default:
		return "", fmt.Errorf("can't convert element (%s) with type (%T) to string : unsupported type conversion", en.path, value)
	}
}

func (en *Entry) AsInt64() (int64, error) {
	v, err := en.AsFloat64()
	return int64(v), err
}

func (en *Entry) AsInt() (int, error) {
	v, err := en.AsFloat64()
	return int(v), err
}

func (en *Entry) AsFloat64() (float64, error) {
	if en.err != nil {
		return 0, fmt.Errorf("can't get element (%s) AsFloat : %w", en.path, en.err)
	}
	switch value := en.entry.(type) {
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case int:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case string:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v, nil
		} else {
			return 0, fmt.Errorf("can't convert string value to float for element (%s) with value (%s) : %w", en.path, value, err)
		}
	default:
		return 0, fmt.Errorf("wrong type for element (%s) : expected int or float but '%T' was found", en.path, value)
	}
}

func (en *Entry) AsBool() (bool, error) {
	if en.err != nil {
		return false, fmt.Errorf("can't get element (%s) AsBool : %w", en.path, en.err)
	}

	switch value := en.entry.(type) {
	case bool:
		return value, nil
	case string:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue, nil
		} else {
			return false, fmt.Errorf("can't parse boolean value for element (%s) with value: (%s)", en.path, value)
		}
	default:
		return false, fmt.Errorf("wrong type for element (%s) : Expected bool or string but %T found", en.path, value)
	}
}

func (en *Entry) AsMap() (map[string]interface{}, error) {
	if en.err != nil {
		return nil, fmt.Errorf("can't get element (%s) AsMap : %w", en.path, en.err)
	}

	if en.entryMap != nil {
		return en.entryMap, nil
	}

	if str, ok := en.entry.(string); ok {
		var mapTemplate interface{}
		err := json.Unmarshal([]byte(str), &mapTemplate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the value of element (%s) to JSON : %w", en.path, err)
		}

		if jsonMap, ok := mapTemplate.(map[string]interface{}); ok {
			en.entryMap = jsonMap
			return jsonMap, nil
		}
		return nil, fmt.Errorf("failed to map JSON values to a map[string]interface{} for element (%s)", en.path)
	}

	if result, ok := en.entry.(map[string]interface{}); ok {
		en.entryMap = result
		return result, nil
	} else {
		return nil, fmt.Errorf("wrong type for element (%s)", en.path)
	}
}

func (en *Entry) AsRawValue() (interface{}, error) {
	if en.err != nil {
		return nil, fmt.Errorf("can't get element (%s) AsRawValue : %w", en.path, en.err)
	}

	return en.entry, nil
}
