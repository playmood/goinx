package gnet

type Message struct {
	Id      uint32 // 消息Id
	DataLen uint32 // 消息的长度
	Data    []byte // 消息内容
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgId 获取消息的Id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetMsgLen 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// GetData 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId 设置消息Id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetData 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
