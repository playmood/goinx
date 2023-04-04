package utils

import (
	"encoding/json"
	"goinx/iface"
	"io/ioutil"
)

// 存储框架全局参数，供其它模块使用
// 某些参数支持goinx.json由用户配置
type GlobalObj struct {
	TcpServer iface.IServer // 全局Server对象
	Host      string        // 主机监听IP
	TcpPort   int           // 监听端口
	Name      string        // 当前服务器名称

	Version        string // goinx版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前框架数据包最大值
}

// 全局对象实例
var GlobalObject *GlobalObj

// 从goinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/goinx.json")
	if err != nil {
		panic(err)
	}
	// 解析json
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 如果配置文件没加载，此为默认值
	GlobalObject = &GlobalObj{
		Name:           "GoinxServerApp",
		Version:        "0.3",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	// 尝试从conf/goinx.json中加载
	GlobalObject.Reload()
}
