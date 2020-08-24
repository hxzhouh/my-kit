/*
@Time : 2020/8/3 11:49
@Author : ZhouHui2
@File : options
@Software: GoLand
*/
package kit

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type AppHandler interface {
	Register(*grpc.Server, *runtime.ServeMux, *AppConfig)
	//Router(*echo.Echo, *AppConfig)

}

type AppConfig struct {
	AppName        string
	Addr           string
	ServerOption   []grpc.ServerOption
	ServeMuxOption []runtime.ServeMuxOption
	AppHandler     AppHandler
}

func NewAppConfig(appName, addr string, appHandler AppHandler) *AppConfig {
	return &AppConfig{
		AppName:    appName,
		Addr:       addr,
		AppHandler: appHandler,
	}
}
