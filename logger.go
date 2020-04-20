/**
* Author: JeffreyBool
* Date: 2020/4/20
* Time: 12:26
* Software: GoLand
 */

package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	ii   *logger
)

type logger struct {
	opts *options
	any  *sync.Map
}

func New(opts ...Option) *logger {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	once.Do(func() {
		ii = &logger{
			opts: &o,
			any:  new(sync.Map),
		}
	})
	return ii
}

func Logger(arg ...string) *zap.Logger {
	if ii == nil {
		panic("logger is not initialized")
	}

	logName := "all"
	if len(arg) > 0 {
		logName = arg[0]
	}

	if log, ok := ii.any.Load(logName); ok {
		return log.(*zap.Logger)
	}

	newLogger, err := newLogger(logName)
	if err != nil || newLogger == nil {
		return nil
	}
	ii.any.Store(logName, newLogger)
	return newLogger
}

func newLogger(filename string) (logger *zap.Logger, err error) {
	if _, err := os.Stat(ii.opts.path); os.IsNotExist(err) {
		if err := os.MkdirAll(ii.opts.path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	var zConf zap.Config
	if ii.opts.env == "online" {
		zConf = NewProductionConfig(filename)
	} else {
		zConf = NewDevelopmentConfig(filename)
	}
	return NewLogger(zConf)
}

func NewProductionConfig(filename string) zap.Config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = encodeTime
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	return zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout", filepath.Join(ii.opts.path, filename) + ".out.log"},
		ErrorOutputPaths:  []string{"stderr", filepath.Join(ii.opts.path, filename) + ".error.log"},
	}
}

func NewDevelopmentConfig(filename string) zap.Config {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = encodeTime
	return zap.Config{
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: true,
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout", filepath.Join(ii.opts.path, filename) + ".out.log"},
		ErrorOutputPaths:  []string{"stderr", filepath.Join(ii.opts.path, filename) + ".error.log"},
	}
}

// 格式化时间
func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
