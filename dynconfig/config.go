package dynconfig

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"sync"
	"time"
)

const DefaultRefreshInterval = time.Minute * 1

//go:generate mockgen -destination mock_test.go -package dynconfig_test . metrics,logger,Adapter

type metrics interface {
	Timing(metricName string, duration time.Duration, tags map[string]string)
	Increment(metricName string, value int, tags map[string]string)
}

type logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Adapter interface {
	GetConfig(key string) (string, error)
}

type scheduler interface {
	RunAfter(time.Duration, func())
}

type record struct {
	value   string
	version string
}

type Config struct {
	adapter         Adapter
	refreshInterval time.Duration
	metrics         metrics
	logger          logger
	scheduler       scheduler
	configMap       map[string]record
	mutex           sync.Mutex
}

type Options struct {
	RefreshInterval time.Duration
	logger          logger
	metrics         metrics
	scheduler       scheduler
}

func NewOptions() *Options {
	return &Options{
		RefreshInterval: DefaultRefreshInterval,
	}
}

func (opts *Options) WithRefreshInterval(refreshInterval time.Duration) *Options {
	opts.RefreshInterval = refreshInterval
	return opts
}

func (opts *Options) WithMetrics(metrics metrics) *Options {
	opts.metrics = metrics
	return opts
}

func (opts *Options) WithLogger(logger logger) *Options {
	opts.logger = logger
	return opts
}

func (opts *Options) WithScheduler(scheduler scheduler) *Options {
	opts.scheduler = scheduler
	return opts
}

func New(configAdapter Adapter, options *Options) *Config {
	logger := options.logger
	if logger == nil {
		logger = &defaultLogger{}
	}

	metrics := options.metrics
	if metrics == nil {
		metrics = &defaultMetrics{}
	}

	scheduler := options.scheduler
	if scheduler == nil {
		scheduler = &defaultScheduler{}
	}

	return new(configAdapter, logger, metrics, scheduler, options)
}

func new(adapter Adapter, logger logger, metrics metrics, scheduler scheduler, options *Options) *Config {
	refreshInterval := DefaultRefreshInterval
	if options != nil && options.RefreshInterval >= time.Second*10 {
		refreshInterval = options.RefreshInterval
	}

	cfg := &Config{
		adapter:         adapter,
		refreshInterval: refreshInterval,
		logger:          logger,
		metrics:         metrics,
		scheduler:       scheduler,
		configMap:       map[string]record{},
		mutex:           sync.Mutex{},
	}

	cfg.scheduleNextRun()
	return cfg
}

func (c *Config) scheduleNextRun() {
	c.scheduler.RunAfter(c.refreshInterval, c.refresh)
}

func (c *Config) getConfigMapCopySync() map[string]record {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	newConfigMap := map[string]record{}
	for k, v := range c.configMap {
		newConfigMap[k] = v
	}
	return newConfigMap
}

func (c *Config) refresh() {
	newConfigMap := c.getConfigMapCopySync()

	for configKey := range newConfigMap {
		value, err := c.fetch(configKey)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn(
					"message", "Failed to fetch a config value",
					"configName", configKey,
					"error", err.Error(),
				)
			}
			continue
		}
		newConfigMap[configKey] = record{
			value:   value,
			version: checksum(value),
		}
	}

	c.swapConfigMapSync(newConfigMap)

	c.scheduleNextRun()
}

func checksum(s string) string {
	return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(s)))
}

func (c *Config) swapConfigMapSync(newConfigMap map[string]record) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k, v := range c.configMap {
		if _, ok := newConfigMap[k]; !ok {
			newConfigMap[k] = v
		}
	}
	c.configMap = newConfigMap
}

func (c *Config) registerSync(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.configMap[key] = record{
		value:   value,
		version: checksum(value),
	}
}

func (c *Config) fetch(key string) (value string, err error) {
	if c.metrics != nil {
		defer func(startTime time.Time) {
			c.metrics.Timing("config_fetch", time.Since(startTime), map[string]string{
				"success": strconv.FormatBool(err == nil),
				"param":   key,
			})
		}(time.Now())
	}

	return c.adapter.GetConfig(key)
}

func (c *Config) get(key string, withMetrics bool) (value string, err error) {
	var version string

	if withMetrics && c.metrics != nil {
		defer func(startTime time.Time) {
			c.metrics.Increment("config_get", 1, map[string]string{
				"success": strconv.FormatBool(err == nil),
				"param":   key,
				"version": version,
			})
		}(time.Now())
	}

	if record, ok := c.getFromConfigMapProtected(key); ok {
		version = record.version
		return record.value, nil
	}

	err = fmt.Errorf("cannot find value for config '%s'", key)
	return "", err
}

func (c *Config) getFromConfigMapProtected(key string) (record, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	record, found := c.configMap[key]
	return record, found
}

func (c *Config) Get(key string) *Entry {
	value, err := c.get(key, true)
	return NewEntryWithRoot(value, key, err)
}

func (c *Config) Fetch(key string) *Entry {
	value, err := c.fetch(key)
	return NewEntryWithRoot(value, key, err)
}

func (c *Config) getAndRegisterStringParam(key string, value string) (string, error) {
	s, err := c.get(key, false)
	if err == nil {
		// The parameter has already been registered
		return s, nil
	}

	s, err = c.adapter.GetConfig(key)
	if err == nil {
		c.registerSync(key, s)
		return s, nil
	}

	c.registerSync(key, value)
	return value, fmt.Errorf("failed to get value for config '%s' : %w", key, err)
}

func (c *Config) GetAndRegisterStringConfig(key string, value string) *Entry {
	s, err := c.getAndRegisterStringParam(key, value)
	return NewEntryWithRoot(s, key, err)
}

func (c *Config) GetAndRegisterBoolConfig(key string, value bool) *Entry {
	return c.GetAndRegisterStringConfig(key, strconv.FormatBool(value))
}

func (c *Config) GetAndRegisterIntConfig(key string, value int) *Entry {
	return c.GetAndRegisterStringConfig(key, strconv.FormatInt(int64(value), 10))
}

func (c *Config) GetAndRegisterFloatConfig(key string, value float64) *Entry {
	return c.GetAndRegisterStringConfig(key, strconv.FormatFloat(value, 'g', -1, 64))
}
