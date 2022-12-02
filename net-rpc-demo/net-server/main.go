package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/sh1yu/go-training/net-rpc-demo/common"
	net_interface "github.com/sh1yu/go-training/net-rpc-demo/net-interface"
)

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "Hello:" + request
	return nil
}

func main() {
	err := net_interface.RegisterHelloService(new(HelloService))
	common.MustNilErr(err, "register service error")

	listener, err := net.Listen("tcp", ":1234")
	common.MustNilErr(err, "ListenTCP error")

	for {
		conn, err := listener.Accept()
		common.MustNilErr(err, "Accept error")

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}
