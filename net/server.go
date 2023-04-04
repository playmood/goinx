package net

import (
	"fmt"
	"goinx/iface"
	"net"
)

// Server IServer接口实例化
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前Server添加一个router，server注册的连接对应的处理业务
	Router iface.IRouter
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listen at IP:%s, Port %d, is starting\n", s.IP, s.Port)

	// 防止阻塞 异步化
	go func() {
		// 1 获取TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err ", err)
			return
		}
		fmt.Println("start Goinx server success, ", s.Name, " listening...")
		var cid uint32 = 0
		// 3 阻塞的等待客户端进行连接， 处理客户端连接业务(读写)
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 将处理新连接的业务方法和conn进行绑定 得到连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 对应连接的业务处理
			go dealConn.Start()
		}
	}()

}

// Stop 停止服务器
func (s *Server) Stop() {
	// todo 将服务器的资源、状态或一些已经开启的连接信息进行回收

}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// todo 做一些启动服务器之后的额外业务

	// 这里阻塞
	select {}
}

func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router
	fmt.Println("add router success!")
}

/*
	初始化Server模块
*/
func NewServer(name string) iface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
}
