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
	err := request.GetConnection().SendMsg(0, []byte("ping .. ping .."))
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

type HelloRouter struct {
	gnet.BaseRouter
}

// Handle test
func (this *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端数据，再回写ping
	fmt.Println("recv from client: msgID: ", request.GetMsgID(), " , data: ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("hello goinx !!"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionBegin 创建连接后执行的钩子函数
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("---> begin hook called ...")
	if err := conn.SendMsg(202, []byte("begin hook called")); err != nil {
		fmt.Println(err)
	}
}

// DoConnectionEnd 销毁连接后执行的钩子函数
func DoConnectionEnd(conn iface.IConnection) {
	fmt.Println("---> end hook called ...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost ...")
}

func main() {
	s := gnet.NewServer("[goinx V0.9]")

	// 注册连接的hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionEnd)

	// 给当前框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
