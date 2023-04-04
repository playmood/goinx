package main

import (
	"fmt"
	"goinx/iface"
	"goinx/net"
)

// PingRouter ping test
type PingRouter struct {
	net.BaseRouter
}

// PreHandle test
func (this *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle test
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ...\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// PostHandle test
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := net.NewServer("[goinx V0.2]")

	// 给当前框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	s.Serve()
}
