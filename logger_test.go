/**
* Author: JeffreyBool
* Date: 2020/4/20
* Time: 12:28
* Software: GoLand
 */

package logger

import (
	"testing"


	"go.uber.org/zap"
)

func TestMain(t *testing.M) {
	New(SetEnv("dev"), SetPath("./log"))
	t.Run()
}

func TestGetLogger(t *testing.T) {
	Logger().Info("msg",zap.String("uid","abc"))
	Logger().Debug("debug",zap.String("uid","abc"))
	Logger().Error("error",zap.String("uid","abc"))
	Logger("goim").Info("info",zap.String("uid","abc"))
	Logger("goim").Error("error",zap.String("uid","abc"))
	Logger("goim").Debug("debug",zap.String("uid","abc"))
}
