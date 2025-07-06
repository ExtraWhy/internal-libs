package dynconfig_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ExtraWhy/internal-libs/dynconfig"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

const testDuration = time.Second * 20

type mockScheduler struct {
	f        func()
	interval time.Duration
	numCalls int
}

func (ms *mockScheduler) RunAfter(interval time.Duration, f func()) {
	ms.numCalls++
	ms.interval = interval
	ms.f = f
}

func (ms *mockScheduler) refresh() {
	ms.f()
}

func TestGetStringParam_NotRegistered_BeforeRefresh(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	mockMetrics.EXPECT().
		Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
			"success": "false",
			"param":   "stringParam",
			"version": "",
		})).Times(1)

	mockAdapter.EXPECT().
		GetConfig(gomock.Eq("stringParam")).Return("", nil).
		Times(0)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	_, err := config.Get("stringParam").AsString()

	assert.NotNil(err)
	assert.Equal("can't get element ([stringParam]) AsString : cannot find value for config 'stringParam'", err.Error())

	assert.Equal(1, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval, "Scheduler should have the correct interval set")
}

func TestGetStringParam_NotRegistered_AfterRefresh(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	mockMetrics.EXPECT().
		Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
			"success": "false",
			"param":   "stringParam",
			"version": "",
		})).Times(1)

	mockAdapter.EXPECT().
		GetConfig(gomock.Eq("stringParam")).Return("", nil).
		Times(0)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	mockScheduler.refresh()

	_, err := config.Get("stringParam").AsString()

	assert.NotNil(err)
	assert.Equal("can't get element ([stringParam]) AsString : cannot find value for config 'stringParam'", err.Error())

	assert.Equal(2, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval, "Scheduler should have the correct interval set")
}

func TestFetchStringParam(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	mockMetrics.EXPECT().
		Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
			"success": "false",
			"param":   "stringParam",
			"version": "",
		})).Times(2)

	mockMetrics.EXPECT().
		Timing(gomock.Eq("config_fetch"), gomock.AssignableToTypeOf(time.Duration(0)), gomock.Any()).Times(2)

	gomock.InOrder(
		mockAdapter.EXPECT().GetConfig(gomock.Eq("stringParam")).Return("value1", nil),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("stringParam")).Return("value2", nil),
	)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	assert.Equal(1, mockScheduler.numCalls)

	mockScheduler.refresh()

	_, err := config.Get("stringParam").AsString()

	assert.NotNil(err)
	assert.Equal("can't get element ([stringParam]) AsString : cannot find value for config 'stringParam'", err.Error())

	assert.Equal(2, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval, "Scheduler should have the correct interval set")

	// Test what happens when we try to get a non-registered string config param from the cache
	value, err := config.Fetch("stringParam").AsString()
	assert.Nil(err)
	assert.Equal("value1", value)

	mockScheduler.refresh()
	assert.Equal(3, mockScheduler.numCalls)

	value, err = config.Fetch("stringParam").AsString()
	assert.Nil(err)
	assert.Equal("value2", value)

	// Validate that the config parameter hasn't been registered or cached
	_, err = config.Get("stringParam").AsString()
	assert.NotNil(err)
	assert.Equal("can't get element ([stringParam]) AsString : cannot find value for config 'stringParam'", err.Error())
}

func TestGetAndRegisterStringParam(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	gomock.InOrder(
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "stringParam",
				"version": "3294772269", // the checksum corresponding to "defaultValue"
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "stringParam",
				"version": "2782488940", // the checksum corresponding to "Value1"
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "stringParam",
				"version": "1020278998", // the checksum corresponding to "Value2"
			})),
	)

	mockMetrics.EXPECT().
		Timing(gomock.Eq("config_fetch"), gomock.AssignableToTypeOf(time.Duration(0)), gomock.Any()).Times(2)

	mockLogger.EXPECT().
		Warn(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gomock.InOrder(
		mockAdapter.EXPECT().GetConfig(gomock.Eq("stringParam")).
			Return("", errors.New("Adapter error")),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("stringParam")).
			Return("Value1", nil),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("stringParam")).
			Return("Value2", nil),
	)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	// Initial fetch with an error - should return the default value after the first error
	_, err := config.GetAndRegisterStringConfig("stringParam", "defaultValue").AsString()
	assert.NotNil(err)
	assert.Equal("can't get element ([stringParam]) AsString : failed to get value for config 'stringParam' : Adapter error", err.Error())

	value, err := config.Get("stringParam").AsString()
	assert.Nil(err)
	assert.Equal("defaultValue", value)

	// This is the first refresh. It should get "Value1" from the remote config through the adapter.
	mockScheduler.refresh()
	value, err = config.Get("stringParam").AsString()
	assert.Nil(err)
	assert.Equal("Value1", value)

	// This is the second refresh. It should get "Value2" from the remote config through the adapter.
	mockScheduler.refresh()
	value, err = config.Get("stringParam").AsString()
	assert.Nil(err)
	assert.Equal("Value2", value)

	assert.Equal(3, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval)
}

func TestGetAndRegisterBoolParam(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	// Then expect success metrics in order
	gomock.InOrder(
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "boolParam",
				"version": "4261170317", // checksum corresponding to "true"
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "boolParam",
				"version": "734881840", // checksum corresponding to "false"
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "boolParam",
				"version": "4261170317", // checksum corresponding to "true"
			})),
	)

	mockMetrics.EXPECT().
		Timing(gomock.Eq("config_fetch"), gomock.AssignableToTypeOf(time.Duration(0)), gomock.Any()).Times(2)

	mockLogger.EXPECT().
		Warn(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gomock.InOrder(
		mockAdapter.EXPECT().GetConfig(gomock.Eq("boolParam")).
			Return("", errors.New("Adapter error")),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("boolParam")).
			Return("false", nil),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("boolParam")).
			Return("true", nil),
	)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	// Initial fetch with error - should return default value after error
	_, err := config.GetAndRegisterBoolConfig("boolParam", true).AsBool()
	assert.NotNil(err)
	assert.Equal("can't get element ([boolParam]) AsBool : failed to get value for config 'boolParam' : Adapter error", err.Error())

	value, err := config.Get("boolParam").AsBool()
	assert.Nil(err)
	assert.Equal(true, value)

	// First refresh - should get "false"
	mockScheduler.refresh()
	value, err = config.Get("boolParam").AsBool()
	assert.Nil(err)
	assert.Equal(false, value)

	// Second refresh - should get "true"
	mockScheduler.refresh()
	value, err = config.Get("boolParam").AsBool()
	assert.Nil(err)
	assert.Equal(true, value)

	assert.Equal(3, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval)
}

func TestGetAndRegisterJSONParam(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockScheduler := &mockScheduler{}
	mockAdapter := NewMockAdapter(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockMetrics := NewMockmetrics(ctrl)

	gomock.InOrder(
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "jsonParam",
				"version": "2745614147", // checksum corresponding to "{}"
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "jsonParam",
				"version": "4089638929", // checksum corresponding to the JSON with key1 and key2
			})),
		mockMetrics.EXPECT().
			Increment(gomock.Eq("config_get"), gomock.Eq(1), gomock.Eq(map[string]string{
				"success": "true",
				"param":   "jsonParam",
				"version": "2745614147", // checksum corresponding to "{}"
			})),
	)

	mockMetrics.EXPECT().
		Timing(gomock.Eq("config_fetch"), gomock.AssignableToTypeOf(time.Duration(0)), gomock.Any()).Times(2)

	mockLogger.EXPECT().
		Warn(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gomock.InOrder(
		mockAdapter.EXPECT().GetConfig(gomock.Eq("jsonParam")).
			Return("", errors.New("Adapter error")),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("jsonParam")).
			Return(`{"key1": "value1", "key2": 2222}`, nil),
		mockAdapter.EXPECT().GetConfig(gomock.Eq("jsonParam")).
			Return("{}", nil),
	)

	options := dynconfig.NewOptions().
		WithRefreshInterval(testDuration).
		WithLogger(mockLogger).
		WithMetrics(mockMetrics).
		WithScheduler(mockScheduler)

	config := dynconfig.New(mockAdapter, options)

	// Initial fetch with error - should return default value after error
	_, err := config.GetAndRegisterStringConfig("jsonParam", "{}").AsMap()
	assert.NotNil(err)
	assert.Equal("can't get element ([jsonParam]) AsMap : failed to get value for config 'jsonParam' : Adapter error", err.Error())

	value, err := config.Get("jsonParam").AsMap()
	assert.Nil(err)
	assert.Equal(map[string]interface{}{}, value)

	// First refresh - should get JSON with key1 and key2
	mockScheduler.refresh()
	value, err = config.Get("jsonParam").AsMap()
	assert.Nil(err)
	assert.Equal(map[string]interface{}{
		"key1": "value1",
		"key2": float64(2222),
	}, value)

	// Second refresh - should get empty JSON object
	mockScheduler.refresh()
	value, err = config.Get("jsonParam").AsMap()
	assert.Nil(err)
	assert.Equal(map[string]interface{}{}, value)

	assert.Equal(3, mockScheduler.numCalls)
	assert.Equal(testDuration, mockScheduler.interval)
}
