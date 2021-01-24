package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"p2pserve/util"
)

type P2PServer struct {
	//WebSocket绑定函数,由信令服务处理
	handleWebSocket func(ws *WebSocketConn, request *http.Request)
	//Websocket升级为长连接
	upgrader websocket.Upgrader
}

//WebSocket请求处理
func (server *P2PServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {
	//返回头
	responseHeader := http.Header{}
	//升级为长连接
	socket, err := server.upgrader.Upgrade(writer, request, responseHeader)

	if err != nil {
		util.Panicf("%v", err)
	}
	//实例化一个WebSocketConn对象
	wsTransport := NewWebSocketConn(socket)
	//处理具体的请求消息
	server.handleWebSocket(wsTransport, request)

	//WebSocketConn开始读取消息
	wsTransport.ReadMessage()

}

//实例化一个P2P服务
func NewP2PServer(wsHandler func(ws *WebSocketConn, request *http.Request)) *P2PServer {
	var server = &P2PServer{
		handleWebSocket: wsHandler,
	}
	server.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return server

}
