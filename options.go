/**
* Author: JeffreyBool
* Date: 2020/4/20
* Time: 12:27
* Software: GoLand
 */

package logger

type Option func(o *options)

type options struct {
	env  string
	path string
}

func SetEnv(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func SetPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}
