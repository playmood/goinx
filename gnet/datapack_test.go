package gnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 测试拆包封包 单元测试
func TestDataPack(t *testing.T) {
	// 创建tcp socket
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("sever listen err: ", err)
		return
	}
	go func() {
		// 从客户端读取数据 进行拆包
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				// 拆包对象
				dp := NewDataPack()
				for {
					// 第一次从conn读，把head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg是有数据的，需要进行第二次读取
						// 第二次读conn，读data len
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						// 完整的一个消息读取完毕
						fmt.Println("--> Recv Msg ID: ", msg.Id, ", datalen = ", msg.DataLen, " data = ", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}
	dp := NewDataPack()
	// 模拟粘包过程，封装两个msg一起发送
	// 第一个包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'g', 'o', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	// 第二个包
	msg2 := &Message{
		Id:      2,
		DataLen: 6,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}
	// 两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 发送
	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("write error")
		return
	}
}
