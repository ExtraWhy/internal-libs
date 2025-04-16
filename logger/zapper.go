package logger

import "go.uber.org/zap"

const (
	NOP      = 0
	DEV      = 1
	PROD     = 2
	INFO     = 1 << 16
	DEBUG    = 1 << 17
	CRITICAL = 1 << 18
)

type ZapperLog struct {
	logger *zap.Logger
}

func (zp *ZapperLog) Init(tp int) {
	switch t := tp; t {
	case NOP: //todo
	case DEV:
		zp.logger = zap.Must(zap.NewDevelopment())
	case PROD:
		zp.logger = zap.Must(zap.NewProduction())
	default:
		zp.logger = zap.Must(zap.NewProduction())
	}
}

func (zp *ZapperLog) Log(level int, message string, zapfields ...zap.Field) {

	defer zp.logger.Sync()
	switch level {
	case INFO:
		zp.logger.Info(message, zapfields...)
	case DEBUG:
		zp.logger.Debug(message, zapfields...)
	case CRITICAL:
		zp.logger.Error(message, zapfields...)
	default:
		zp.logger.Debug(message, zapfields...)
	}
}
