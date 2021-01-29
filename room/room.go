package room

import (
	"net/http"
	"p2pserve/messagetype"
	"p2pserve/server"
	"p2pserve/util"
)

//定义房间
type Room struct {
	//所有用户
	users map[string]User
	//所有会话
	sessions map[string]Session
	//房间ID
	ID string
}

func NewRoom(id string) *Room {
	var room = &Room{
		users:    make(map[string]User),
		sessions: make(map[string]Session),
		ID:       id,
	}
	return room
}

//定义房间管理
type RoomManager struct {
	rooms map[string]*Room
}

//获取房间
func (roomManager *RoomManager) getRoom(id string) *Room {
	return roomManager.rooms[id]
}

//创建房间
func (roomManager *RoomManager) createRoom(id string) *Room {
	roomManager.rooms[id] = NewRoom(id)
	return roomManager.rooms[id]
}

//删除房间
func (roomManager *RoomManager) deleteRoom(id string) {
	delete(roomManager.rooms, id)
}

//WebSocket消息处理
func (roomManager *RoomManager) HandleNewWebSocket(conn *server.WebSocketConn, request *http.Request) {
	util.Infof("On Open %v", request)

	//连接关闭事件处理
	conn.On("close", func(code int, text string) {
		onClose(conn, roomManager)
	})

}

//通知所有的用户更新
func (roomManager *RoomManager) notifyUserUpdate(conn *server.WebSocketConn, users map[string]User) {
	infos := []UserInfo{}
	for _, userClient := range users {
		infos = append(infos, userClient.info)
	}
	request := make(map[string]interface{})
	request["type"] = messagetype.UpdateUserList.String()
	//数据
	request["data"] = infos
	//迭代所有的User
	for _, user := range users {
		//将Json数据发送给每一个User
		user.conn.Send(util.Marshal(request))
	}
}

//连接关闭
func onClose(conn *server.WebSocketConn, roomManager *RoomManager) {
	util.Infof("连接关闭 %v", conn)
	var userId string = ""
	var roomId string = ""

	//遍历所有的房间找到退出的用户
	for _, room := range roomManager.rooms {
		for _, user := range room.users {
			if user.conn == conn {
				userId = user.info.ID
				roomId = room.ID
			}
		}
	}
	//end 遍历所有的房间找到退出的用户
	if roomId == "" {
		util.Errorf("没有查找到退出的房间及用户")
		return
	}
	util.Infof("退出的用户roomId %v userId %v", roomId, userId)

	//循环遍历所有的User
	for _, user := range roomManager.getRoom(roomId).users {
		if user.conn != conn {
			level := map[string]interface{}{
				"type": messagetype.LeaveRoom.String(),
				"data": userId,
			}
			user.conn.Send(util.Marshal(level))
		}
	}
	util.Infof("User退出", userId)

	//根据Id删除User
	delete(roomManager.getRoom(roomId).users, userId)

	//通知所有的User更新数据

	roomManager.notifyUserUpdate(conn, roomManager.getRoom(roomId).users)

}
