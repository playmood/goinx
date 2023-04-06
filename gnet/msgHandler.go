package gnet

import (
	"fmt"
	"goinx/iface"
	"goinx/utils"
	"strconv"
)

// MsgHandle 消息处理实现层
type MsgHandle struct {
	// 存放msg id 对应的处理方法
	ApiSet map[uint32]iface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	// 工作池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandle 构造函数
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		ApiSet: make(map[uint32]iface.IRouter),
		// 从全局配置获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度/执行Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	// 从request中找到msg id
	handler, ok := mh.ApiSet[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found, please register")
	}
	// 根据msg id 调用对应的Router任务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router iface.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.ApiSet[msgID]; ok {
		panic("double register api: " + strconv.Itoa(int(msgID)))
	}
	// 添加msg与API的绑定关系
	mh.ApiSet[msgID] = router
	fmt.Println("add api msg id = ", msgID, " successful")
}

// StartWorkerPool 启动一个worker工作池，单次初始化
func (mh *MsgHandle) StartWorkerPool() {
	// 根据size分别开启worker，每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 当前worker的消息队列
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxTaskSize)
		// 启动当前worker，阻塞等待消息从channel传入
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iface.IRequest) {
	fmt.Println("worker id = ", workerID, "is started...")
	// 阻塞等待
	for {
		select {
		// 如果有消息过来，执行request绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request iface.IRequest) {
	// 将消息平均分配给worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("add ConnID = ", request.GetConnection().GetConnID(), " request msg id = ",
		request.GetMsgID(), " to Worker ID = ", workerID)

	// 将消息发送给对应worker的taskQueue即可
	mh.TaskQueue[workerID] <- request
}
