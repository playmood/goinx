package gnet

import (
	"fmt"
	"goinx/iface"
	"strconv"
)

// MsgHandle 消息处理实现层
type MsgHandle struct {
	// 存放msg id 对应的处理方法
	ApiSet map[uint32]iface.IRouter
}

// NewMsgHandle 构造函数
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		ApiSet: make(map[uint32]iface.IRouter),
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
