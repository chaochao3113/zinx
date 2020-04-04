/**
* @Author: Chao
* @Date: 2020-03-12 20:27
* @Version: 1.0
 */

package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于Zinx框架来开发的服务器端应用程序
 */

//ping test自定义路由
type PingRouter struct {
	znet.BaseRouter
}


// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest)  {
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据,再回写ping...ping...ping
	fmt.Println("recv from client : msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}


func main() {
	// 1 创建一个server的句柄,使用Zinx的api
	//s := znet.NewServer("[zinx V0.3]")
	s := znet.NewServer()
	// 2 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3 启动server
	s.Serve()
}
