package iface

// 请求模块的抽象层
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
