package iface

// 消息管理抽象层
type IMsgHandle interface {
	// DoMsgHandler 调度/执行Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
