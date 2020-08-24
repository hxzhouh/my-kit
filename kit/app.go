/*
@Time : 2020/8/3 11:30
@Author : ZhouHui2
@File : app
@Software: GoLand
*/
package kit

import "google.golang.org/grpc"

type app struct {
	rpcServer *grpc.Server // grpc 服务端
}

func NewApp() *app {
	return &app{}
}

func (t *app) Init() error {
	t.rpcServer = grpc.NewServer() // 创建一个新的grpc服务
	return nil
}
