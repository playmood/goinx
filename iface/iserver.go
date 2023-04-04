package iface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动
	Start()
	// Stop 停止
	Stop()
	// Server 运行服务器
	Serve()
}
