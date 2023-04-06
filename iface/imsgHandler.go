package iface

// 消息管理抽象层
type IMsgHandle interface {
	// DoMsgHandler 调度/执行Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动工作池
	StartWorkerPool()
	// StartOneWorker 启动一个worker
	StartOneWorker(workerID int, taskQueue chan IRequest)
	// SendMsgToTaskQueue 将消息发给交给消息队列
	SendMsgToTaskQueue(request IRequest)
}
