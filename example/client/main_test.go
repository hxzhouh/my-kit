/*
@Time : 2020/8/24 17:42
@Author : ZhouHui2
@File : main_test.go
@Software: GoLand
*/
package client

import (
	"context"
	"fmt"
	"github.com/hxzhouh/my-kit/example/protobuf"
	"github.com/hxzhouh/my-kit/kit/discovery/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	schema, err := resolver.GenerateAndRegisterConsulResolver("192.168.1.117:8500", "HelloService")
	if err != nil {
		log.Fatal("init consul resovler err", err.Error())
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:///HelloService", schema), grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := protobuf.NewHelloClient(conn)

	// Contact the server and print out its response.
	name := "nosixtools"

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.SayHello(ctx, &protobuf.HelloReq{
			Name: name})
		if err != nil {
			log.Println(fmt.Sprintf("could not greet: %v", err))

		} else {
			log.Printf("Hello: %s", r.Msg)
		}
		time.Sleep(time.Second)
	}
}
