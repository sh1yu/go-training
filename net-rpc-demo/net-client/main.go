package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/sh1yu/go-training/net-rpc-demo/common"
	net_interface "github.com/sh1yu/go-training/net-rpc-demo/net-interface"
)

func main() {
	// 以下好像会出现unexpected EOF
	//c, err := net_interface.DialHelloService("tcp", "localhost:1234")
	//common.MustNilErr(err, "dialing error")
	//var reply string
	//err = c.Hello("hello", &reply)
	//common.MustNilErr(err, "")

	conn, err := net.Dial("tcp", ":1234")
	common.MustNilErr(err, "dialing error")
	c := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	err = c.Call(net_interface.HelloServiceName+".Hello", "hello", &reply)
	common.MustNilErr(err, "")

	fmt.Println(reply)
}
