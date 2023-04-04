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
	AddRouter(router IRouter)
}
