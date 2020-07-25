package logger

// Option 日志配置选项
type Option func(o *options)

type options struct {
	env    string
	path   string
	output bool
	debug  bool
}

// SetEnv 设置日志环境变量
func SetEnv(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

// SetPath 设置日志路径
func SetPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}

// SetOutput 设置是否标准输出
func SetOutput(output bool) Option {
	return func(o *options) {
		o.output = output
	}
}

// SetDebug 设置日志是否开启调试模式
func SetDebug(debug bool) Option {
	return func(o *options) {
		o.debug = debug
	}
}
