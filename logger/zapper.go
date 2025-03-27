package logger

import "go.uber.org/zap"

const (
	NOP  = 0
	DEV  = 1
	PROD = 2
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
		zp.logger = zap.Must(zap.NewDevelopment())
	}
}

func (zp *ZapperLog) Deinit() {
	zp.logger.Sync()
}

func (zp *ZapperLog) GetLog() *zap.Logger {
	return zp.logger //just use zapper log featurse
}
