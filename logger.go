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

// New 初始化日志实例
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

// Logger 获取日志实例
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
	if ii.opts.debug {
		zConf = newDevelopmentConfig(filename)
	} else {
		zConf = newProductionConfig(filename)
	}
	return NewLogger(zConf)
}

func newProductionConfig(filename string) zap.Config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = encodeTime
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	var outputPaths, errorOutputPaths []string
	if ii.opts.output {
		outputPaths = []string{"stdout", filepath.Join(ii.opts.path, filename) + ".out.log"}
		errorOutputPaths = []string{"stderr", filepath.Join(ii.opts.path, filename) + ".error.log"}
	} else {
		outputPaths = []string{filepath.Join(ii.opts.path, filename) + ".out.log"}
		errorOutputPaths = []string{filepath.Join(ii.opts.path, filename) + ".error.log"}
	}
	return zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
	}
}

func newDevelopmentConfig(filename string) zap.Config {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = encodeTime
	var outputPaths, errorOutputPaths []string
	if ii.opts.output {
		outputPaths = []string{"stdout", filepath.Join(ii.opts.path, filename) + ".out.log"}
		errorOutputPaths = []string{"stderr", filepath.Join(ii.opts.path, filename) + ".error.log"}
	} else {
		outputPaths = []string{filepath.Join(ii.opts.path, filename) + ".out.log"}
		errorOutputPaths = []string{filepath.Join(ii.opts.path, filename) + ".error.log"}
	}
	return zap.Config{
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: true,
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
	}
}

// 格式化时间
func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
