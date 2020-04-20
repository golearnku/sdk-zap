/**
* Author: JeffreyBool
* Date: 2020/4/20
* Time: 12:28
* Software: GoLand
 */

package logger_test

import (
	"testing"

	logger "github.com/golearnku/sdk-zap"
	"go.uber.org/zap"
)

func TestMain(t *testing.M) {
	logger.New(logger.SetEnv("dev"), logger.SetPath("./log"))
	t.Run()
}

func TestGetLogger(t *testing.T) {
	logger.Logger().Info("msg", zap.String("uid", "abc"))
	logger.Logger().Debug("debug", zap.String("uid", "abc"))
	logger.Logger().Error("error", zap.String("uid", "abc"))

	// 多实例日志
	logger.Logger("goim").Info("info", zap.String("uid", "abc"))
	logger.Logger("goim").Error("error", zap.String("uid", "abc"))
	logger.Logger("goim").Debug("debug", zap.String("uid", "abc"))
}
