package net

import (
	"fmt"
	"goinx/iface"
	"io"
	"net"
)

// Connection 连接模块
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的id
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 告知当前连接已经退出的channel
	ExitChan chan bool
	// 改连接处理的方法Router
	Router iface.IRouter
}

// NewConnection 初始化连接模块
func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")
	defer fmt.Println("connId = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Println("recv buf err", err)
			continue
		}
		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 执行注册的路由方法
		go func(request iface.IRequest) {
			// 从路由中，找到注册绑定的conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("conn start... connID = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	// todo 启动从当前连接写数据的业务

}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("conn stop... connID = ", c.ConnID)

	// 如果当前连接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP 状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
