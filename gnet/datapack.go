package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"goinx/iface"
	"goinx/utils"
)

// DataPack 封包、拆包的具体模块
type DataPack struct {
}

// NewDataPack 初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取报的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 4字节 + ID uint32 4字节
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	// 创建存放字节的缓冲区
	dataBuffer := bytes.NewBuffer([]byte{})
	// 将dataLen写入databuf中
	err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}
	// 写MsgId
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	// 写数据
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// Unpack 拆包方法
func (dp *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 解压head信息，得到data len和msg id
	msg := &Message{}

	// 读message len
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msg id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断data len 是否超出了允许的最大长度
	if msg.DataLen > utils.GlobalObject.MaxPackageSize && utils.GlobalObject.MaxPackageSize > 0 {
		return nil, errors.New("message is too large")
	}

	return msg, nil
}
