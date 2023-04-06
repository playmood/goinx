package iface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动
	Start()
	// Stop 停止
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 给当前的服务注册一个路由方法，供客户端连接调用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前Server的conn管理器
	GetConnMgr() IConnManager
	// SetOnConnStart 注册OnConnStart钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	// SetOnConnStop 注册OnConnStop钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStart 调用OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用OnConnStop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
