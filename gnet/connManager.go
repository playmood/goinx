package gnet

import (
	"errors"
	"fmt"
	"goinx/iface"
	"sync"
)

// ConnManager 连接管理模块
type ConnManager struct {
	// 管理的连接集合
	Connections map[uint32]iface.IConnection
	// 保护连接集合的读写锁
	ConnLock sync.RWMutex
}

// NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]iface.IConnection),
	}
}

// Add 添加连接
func (connMgr *ConnManager) Add(conn iface.IConnection) {
	// 保护共享资源，加写锁
	connMgr.ConnLock.Lock()
	defer connMgr.ConnLock.Unlock()

	connMgr.Connections[conn.GetConnID()] = conn
	fmt.Println("======>> connID = ", conn.GetConnID(), " conn add to connMgr successfully: conn num = ", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn iface.IConnection) {
	// 加写锁
	connMgr.ConnLock.Lock()
	defer connMgr.ConnLock.Unlock()

	// 删除连接信息
	delete(connMgr.Connections, conn.GetConnID())
	fmt.Println("======>> connID = ", conn.GetConnID(), " remove from connMgr successfully: conn num = ", connMgr.Len())
}

// Get 根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	// 保护共享资源map，加读锁
	connMgr.ConnLock.RLock()
	defer connMgr.ConnLock.RUnlock()

	if conn, ok := connMgr.Connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn not FOUND")
	}
}

// Len 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.Connections)
}

// ClearConn 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 加写锁
	connMgr.ConnLock.RLock()
	defer connMgr.ConnLock.RUnlock()

	for connID, conn := range connMgr.Connections {
		conn.Stop()
		delete(connMgr.Connections, connID)
	}

	fmt.Println("clear all conn successful, num = ", connMgr.Len())
}
