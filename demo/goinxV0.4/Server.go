package main

import (
	"fmt"
	"goinx/gnet"
	"goinx/iface"
)

// PingRouter ping test
type PingRouter struct {
	gnet.BaseRouter
}

// PreHandle test
func (this *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	//if err != nil {
	//	fmt.Println("call back before ping error")
	//}
}

// Handle test
func (this *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端数据，再回写ping
	fmt.Println("recv from client: msgID: ", request.GetMsgID(), " , data: ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping .. ping .."))
	if err != nil {
		fmt.Println(err)
	}
}

// PostHandle test
func (this *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	//if err != nil {
	//	fmt.Println("call back after ping error")
	//}
}

func main() {
	s := gnet.NewServer("[goinx V0.4]")

	// 给当前框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	s.Serve()
}
