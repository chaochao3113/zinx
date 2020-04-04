/**
* @Author: Chao
* @Date: 2020-03-12 20:27
* @Version: 1.0
 */

package main

import "zinx/znet"

/*
	基于Zinx框架来开发的服务器端应用程序
 */

func main() {
	// 1 创建一个server的句柄,使用Zinx的api
	//s := znet.NewServer("[zinx V0.2]")
	s := znet.NewServer()
	// 2 启动server
	s.Serve()
}
