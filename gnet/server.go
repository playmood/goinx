package gnet

import (
	"fmt"
	"goinx/iface"
	"goinx/utils"
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
	// 当前server的消息管理模块，统一管理Router
	MsgHandler iface.IMsgHandle
	// 该server的连接管理模块
	ConnMgr iface.IConnManager
	// 创建连接时自动调用的hook函数
	OnConnStart func(iface.IConnection)
	// 销毁连接时自动调用的hook函数
	OnConnStop func(iface.IConnection)
}

/*
	初始化Server模块
*/
func NewServer(name string) iface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
}

// SetOnConnStart 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection iface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection iface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> call CallOnConnStart() ...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----> call CallOnConnStop() ...")
		s.OnConnStop(conn)
	}
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Goinx Start] Server Name: %s, Listen at IP:%s, Port %d, is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Goinx] Version: %s, MaxConn: %d, MaxSize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	// 防止阻塞 异步化
	go func() {
		// 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()
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

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// todo 给客户端响应一个超出最大连接的错误包
				fmt.Println("======>> too many conn, Maxconn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定 得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 对应连接的业务处理
			go dealConn.Start()
		}
	}()

}

// Stop 停止服务器
func (s *Server) Stop() {
	// 将服务器的资源、状态或一些已经开启的连接信息进行回收
	fmt.Println("[stop] Goinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// todo 做一些启动服务器之后的额外业务

	// 这里阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("add router success!")
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}
