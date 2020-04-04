/**
* @Author: Chao
* @Date: 2020-03-21 15:56
* @Version: 1.0
 */

package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负载Worker读取任务的消息队列
	TaskQuene []chan ziface.IRequest

	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQuene:make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! Need Register!")
		return
	}
	//2 根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		fmt.Println("repeat api, msgID = ", msgID)
		return
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add API MsgID = ", msgID, "succ!")
}

//启动一个worker工作池(开启工作池的动作只能发生一次,一个zinx框架只能有一个worker工作池)
func (mh *MsgHandle) StartWorkerPoll() {
	//根据workerPoolSize分别开启Worker,每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前的worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel...
		mh.TaskQuene[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker,阻塞等待消息从channel传递进来
		go mh.startOneWorker(i, mh.TaskQuene[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	//不断的阻塞等待对应的消息队列的消息
	for {
		select {
		//如果有消息过来,出列的就是一个客户端的Request,执行当前Request所绑定的业务
		case request := <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue (request ziface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的ConnID进行分配
	workID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(), " to WorkerID = ", workID)

	//2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQuene[workID] <- request
}