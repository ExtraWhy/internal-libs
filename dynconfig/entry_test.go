package dynconfig_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ExtraWhy/internal-libs/dynconfig"
	"github.com/stretchr/testify/assert"
)

func TestConfigEntry_Error(t *testing.T) {
	assert := assert.New(t)

	entry := dynconfig.NewEntryWithRoot("value", "root", errors.New("test error"))
	_, err := entry.AsString()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]) AsString : test error", err.Error())

	entry = entry.Get("nested")
	_, err = entry.AsString()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]/nested) AsString : test error", err.Error())

	entry = dynconfig.NewEntryWithRoot("value", "root", nil).Get("nested")
	_, err = entry.AsString()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]/nested) AsString : element not found/reachable (not a JSON/Map)", err.Error())
}

func TestConfigEntry_StringValue(t *testing.T) {
	testSetup := []struct {
		value         interface{}
		expectedValue interface{}
		expectedError error
	}{
		{"value", "value", nil},
		{"", "", nil},
		{true, "true", nil},
		{false, "false", nil},
		{125, "125", nil},
		{-125, "-125", nil},
		{125.5, "125.5", nil},
		{-125.5, "-125.5", nil},
		{map[string]interface{}{"key1": "value1"}, nil, errors.New("can't convert element ([root]) with type (map[string]interface {}) to string : unsupported type conversion")},
	}

	for _, setup := range testSetup {
		t.Run(fmt.Sprintf("TestConfigEntry_StringValue_%v", setup.value), func(t *testing.T) {
			assert := assert.New(t)
			entry := dynconfig.NewEntryWithRoot(setup.value, "root", nil)
			stringValue, err := entry.AsString()
			if setup.expectedError != nil {
				assert.NotNil(err)
				assert.Equal(setup.expectedError.Error(), err.Error())
			} else {
				assert.Nil(err)
				assert.Equal(setup.expectedValue, stringValue)
			}
		})
	}
}

func TestConfigEntry_BoolValue(t *testing.T) {
	testSetup := []struct {
		value         interface{}
		expectedValue interface{}
		expectedError error
	}{
		{"FaLsE", false, errors.New("can't parse boolean value for element ([root]) with value: (FaLsE)")},
		{"false", false, nil},
		{"False", false, nil},
		{"FALSE", false, nil},
		{"TruE", false, errors.New("can't parse boolean value for element ([root]) with value: (TruE)")},
		{"true", true, nil},
		{"True", true, nil},
		{"TRUE", true, nil},
		{true, true, nil},
		{false, false, nil},
	}

	for _, setup := range testSetup {
		t.Run(fmt.Sprintf("TestConfigEntry_BoolValue_%v", setup.value), func(t *testing.T) {
			assert := assert.New(t)
			entry := dynconfig.NewEntryWithRoot(setup.value, "root", nil)
			boolValue, err := entry.AsBool()
			if setup.expectedError != nil {
				assert.NotNil(err)
				assert.Equal(setup.expectedError.Error(), err.Error())
			} else {
				assert.Nil(err)
				assert.Equal(setup.expectedValue, boolValue)
			}
		})
	}
}

func TestEntry_FloatValue(t *testing.T) {
	testSetup := []struct {
		value         interface{}
		expectedValue interface{}
		expectedError error
	}{
		{"125.5", 125.5, nil},      // possitive float
		{-125.5, -125.5, nil},      // negative float
		{"125", float64(125), nil}, // int to float
		{125, float64(125), nil},   // string int representation to float
		{"1E3", float64(1000), nil},
		{"invalid", nil, errors.New(`can't convert string value to float for element ([root]) with value (invalid) : strconv.ParseFloat: parsing "invalid": invalid syntax`)},
		{map[string]interface{}{"key1": "value1"}, nil, errors.New(`wrong type for element ([root]) : expected int or float but 'map[string]interface {}' was found`)},
		{true, nil, errors.New(`wrong type for element ([root]) : expected int or float but 'bool' was found`)},
	}

	for _, setup := range testSetup {
		t.Run(fmt.Sprintf("TestConfigEntry_FloatValue_%v", setup.value), func(t *testing.T) {
			assert := assert.New(t)
			entry := dynconfig.NewEntryWithRoot(setup.value, "root", nil)
			floatValue, err := entry.AsFloat64()
			if setup.expectedError != nil {
				assert.NotNil(err)
				assert.Equal(setup.expectedError.Error(), err.Error())
			} else {
				assert.Nil(err)
				assert.Equal(setup.expectedValue, floatValue)
			}
		})
	}
}

func TestEntry_IntValue(t *testing.T) {
	testSetup := []struct {
		value         interface{}
		expectedValue interface{}
		expectedError error
	}{
		{"125.5", int64(125), nil},
		{-125.5, int64(-125), nil},
		{"125", int64(125), nil},
		{125, int64(125), nil},
		{"1E3", int64(1000), nil},
		{"invalid", nil, errors.New(`can't convert string value to float for element ([root]) with value (invalid) : strconv.ParseFloat: parsing "invalid": invalid syntax`)},
		{map[string]interface{}{"key1": "value1"}, nil, errors.New(`wrong type for element ([root]) : expected int or float but 'map[string]interface {}' was found`)},
		{true, nil, errors.New(`wrong type for element ([root]) : expected int or float but 'bool' was found`)},
	}

	for _, setup := range testSetup {
		t.Run(fmt.Sprintf("TestConfigEntry_IntValue_%v", setup.value), func(t *testing.T) {
			assert := assert.New(t)
			entry := dynconfig.NewEntryWithRoot(setup.value, "root", nil)
			intValue, err := entry.AsInt64()
			if setup.expectedError != nil {
				assert.NotNil(err)
				assert.Equal(setup.expectedError.Error(), err.Error())
			} else {
				assert.Nil(err)
				assert.Equal(setup.expectedValue, intValue)
			}
		})
	}
}

func TestEntry_StringToMap(t *testing.T) {
	assert := assert.New(t)

	var obj interface{}
	json.Unmarshal([]byte("{\"strKey\": \"strValue\", \"intKey\": 22}"), &obj)

	m, err := dynconfig.NewEntryWithRoot("{}", "root", nil).AsMap()
	assert.Nil(err)
	assert.Equal(map[string]interface{}{}, m)

	m, err = dynconfig.NewEntryWithRoot("{\"strKey\": \"strValue\", \"intKey\": 22}", "root", nil).AsMap()
	assert.Nil(err)
	assert.Equal(map[string]interface{}{"strKey": "strValue", "intKey": float64(22)}, m)
}

func TestEntry(t *testing.T) {
	testMap := map[string]interface{}{
		"string": "value1",
		"number": 5005,
		"float":  2002.2,
		"bool":   true,
		"nested": map[string]interface{}{
			"string": "value2",
			"number": 3003,
			"bool":   false,
		},
	}

	assert := assert.New(t)

	strValue, err := dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("string").AsString()
	assert.Nil(err)
	assert.Equal("value1", strValue)

	intValue, err := dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("number").AsInt()
	assert.Nil(err)
	assert.Equal(5005, intValue)
	int64Value, err := dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("number").AsInt64()
	assert.Nil(err)
	assert.Equal(int64(5005), int64Value)

	boolValue, err := dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("bool").AsBool()
	assert.Nil(err)
	assert.Equal(true, boolValue)

	floatValue, err := dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("float").AsFloat64()
	assert.Nil(err)
	assert.Equal(2002.2, floatValue)
	intValue, err = dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("float").AsInt()
	assert.Nil(err)
	assert.Equal(2002, intValue)

	// Verify getting a non-existing key
	_, err = dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("non_existing").AsBool()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]/non_existing) AsBool : element not found", err.Error())

	// Verify getting a non-existing path with bigger depth
	_, err = dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("test1").Get("test2").Get("test3").AsBool()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]/test1/test2/test3) AsBool : element not found", err.Error())

	// Verify getting a nested key out of a simple-type key results in an error
	_, err = dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("bool").Get("test2").Get("test3").AsBool()
	assert.NotNil(err)
	assert.Equal("can't get element ([root]/bool/test2/test3) AsBool : element not found/reachable (not a JSON/Map)", err.Error())

	intValue, err = dynconfig.NewEntryWithRoot(testMap, "root", nil).Get("nested").Get("number").AsInt()
	assert.Nil(err)
	assert.Equal(3003, intValue)
}

func TestConfig(t *testing.T) {
	const testConfig = `{
		"key1": {
		  "stringProp": "stringValue",
		  "intProp": 570
		}
	  }`

	assert := assert.New(t)

	ent := dynconfig.NewEntryWithRoot(testConfig, "/", nil).Get("key1")
	assert.True(ent.Exists())

	providerLimit, err := ent.Get("intProp").AsInt()
	assert.Nil(err)
	assert.Equal(570, providerLimit)
}
