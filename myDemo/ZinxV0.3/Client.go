/**
* @Author: Chao
* @Date: 2020-03-13 14:52
* @Version: 1.0
 */

package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
 */

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	// 1 直接连接远程服务器,得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 2 连接调用Write 写数据
		_, err := conn.Write([]byte("Hello Zinx V0.3..."))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
