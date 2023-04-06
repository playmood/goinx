package gnet

import (
	"errors"
	"fmt"
	"goinx/iface"
	"goinx/utils"
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
	// 无缓冲管道用于读写goroutine同步
	msgChan chan []byte
	// 告知当前连接已经退出的channel，由Reader告知Writer
	ExitChan chan bool
	// 统一管理router
	MsgHandler iface.IMsgHandle
}

// NewConnection 初始化连接模块
func NewConnection(conn *net.TCPConn, connID uint32, handler iface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
}

// StartWriter 写消息的goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())
	// 不断阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error, ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，Writer也要退出
			return
		}
	}
}

// StartReader 读消息的goroutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("[conn Reader exit!] remote addr is ", c.RemoteAddr().String(), " connId = ", c.ConnID)
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	if err == io.EOF {
		//		continue
		//	}
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		// 创建包对象
		dp := NewDataPack()
		// 读取客户端的msg head 8B
		headData := make([]byte, dp.GetHeadLen())
		// 拆包，得到msgID 和 msgDatalen 放在msg消息中
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error")
			break
		}
		msg, err := dp.Unpack(headData)
		if msg == nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据data len， 再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 将消息发给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 根据绑定好的msg id 找到 对应的api router执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("conn start... connID = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()

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

	// 告知Writer关闭
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// SendMsg 提供一个SendMsg方法 将发送给客户端的数据，先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("conn closed when send msg")
	}
	// 将data进行封包 MsgDataLen MsgID Data
	dp := NewDataPack()
	binMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack data failed!")
		return errors.New("pack data failed")
	}
	// 将数据发送给管道
	c.msgChan <- binMsg

	return nil
}
