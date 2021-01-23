package server

import (
	"github.com/chuckpreslar/emission"
	"github.com/gorilla/websocket"
	"p2pserve/util"
	"sync"
	"time"
)

//发送心跳包的间隔时间 5秒
const pingPeriod = 5 * time.Second

type WebSocketConn struct {
	//事件派发器
	emission.Emitter
	//socket连接
	socket *websocket.Conn
	//互斥锁
	mutex *sync.Mutex
	//是否关闭
	closed bool
}

//实例化WebSocket连接
func NewWebSocketConn(socket *websocket.Conn) *WebSocketConn {
	var conn WebSocketConn
	conn.Emitter = *emission.NewEmitter()
	//socket连接
	conn.socket = socket
	//实例化互斥锁
	conn.mutex = new(sync.Mutex)
	//打开状态
	conn.closed = false

	//socket连接关闭回调函数
	conn.socket.SetCloseHandler(func(code int, text string) error {
		//输出日志
		util.Warnf("%s [%d]", text, code)
		//派发关闭事件
		conn.Emit("close", code, text)
		//设置为关闭状态
		conn.closed = true
		return nil
	})
	//返回连接
	return &conn
}
