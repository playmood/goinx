package main

import (
	"fmt"
	"net"
	"time"
)

/*
模拟客户端连接
*/

func main() {
	fmt.Println("client start")
	time.Sleep(1 * time.Second)
	// 1 直接连接远程服务器 得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error, exit!")
	}
	// 2 连接调用write写数据
	for {
		_, err := conn.Write([]byte("hello goinx"))
		if err != nil {
			fmt.Println("write conn error", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("server answer: %s, cnt = %d \n", buf, cnt)

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
