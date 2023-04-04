package iface

// 路由抽象层
type IRouter interface {
	// 处理业务之前的方法hook
	PreHandle(request IRequest)
	// 处理业务的主方法hook
	Handle(request IRequest)
	// 处理完conn业务之后的方法hook
	PostHandle(request IRequest)
}
