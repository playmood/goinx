package main

import (
	"fmt"
	"goinx/gnet"
	"io"
	"net"
	"time"
)

/*
模拟客户端连接
*/

func main() {
	fmt.Println("client1 start")
	time.Sleep(3 * time.Second)
	// 1 直接连接远程服务器 得到一个conn连接
	conn, err := net.Dial("tcp", "192.168.18.130:8999")
	if err != nil {
		fmt.Println("client start error, exit!")
	}
	// 2 连接调用write写数据
	for {
		// 发送封包的message
		dp := gnet.NewDataPack()
		msg, err := dp.Pack(gnet.NewMessage(1, []byte("I want hello")))
		if err != nil {
			fmt.Println("pack error", err)
			return
		}
		//fmt.Println(string(msg))
		//msg1, err := dp.Unpack(msg)
		//fmt.Println(msg1.GetMsgId(), msg1.GetMsgLen())
		if _, err := conn.Write(msg); err != nil {
			fmt.Println("write error", err)
			return
		}
		// 服务器回复message, 1, ping ping
		// 先读取流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error", err)
			break
		}
		// 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		// 再根据data len进行第二次读取，将data读取出来
		if msgHead.GetMsgLen() > 0 {
			// msg有数据
			msg := msgHead.(*gnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error, ", err)
				return
			}
			fmt.Println("--> recv server Msg: Id = ", msg.Id, " , len = ", msg.DataLen, ", data = ", string(msg.Data))
		}
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
