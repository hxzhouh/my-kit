/*
@Time : 2020/8/24 17:09
@Author : ZhouHui2
@File : main
@Software: GoLand
*/
package main

import (
	"context"
	"fmt"
	"github.com/hxzhouh/my-kit/kit/discovery/register"

	"github.com/hxzhouh/my-kit/example/protobuf"
	"github.com/hxzhouh/my-kit/kit/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	host        = "10.0.109.207"
	consul_host = "192.168.1.117"
	port        = 8081
	consul_port = 8500
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *protobuf.HelloReq) (*protobuf.HelloResp, error) {
	fmt.Println("client called! 8081")
	return &protobuf.HelloResp{Msg: "hi," + in.Name + "!"}, nil
}
func main() {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(host), port, ""})
	if err != nil {
		fmt.Println(err.Error())
	}
	s := grpc.NewServer()

	// register service
	cr := register.NewConsulRegister(fmt.Sprintf("%s:%d", consul_host, consul_port), 15)
	cr.Register(discovery.RegisterInfo{
		Host:           host,
		Port:           port,
		ServiceName:    "HelloService",
		UpdateInterval: time.Second})

	protobuf.RegisterHelloServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		fmt.Println("failed to serve:" + err.Error())
	}
}
