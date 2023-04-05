package iface

// 将请求的消息封装到一个message中，定义抽象的接口
type IMessage interface {
	GetMsgId() uint32  // 获取消息的Id
	GetMsgLen() uint32 // 获取消息长度
	GetData() []byte   // 获取消息内容
	SetMsgId(uint32)   //设置消息Id
	SetData([]byte)    // 设置消息内容
	SetDataLen(uint32) // 设置消息的长度
}
