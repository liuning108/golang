package server

import (
	"fmt"
	"net/http"
	"p2pserve/messagetype"
	"p2pserve/util"
	"strings"
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

//实例化房间管理对象
func NewRoomManager() *RoomManager {
	var roomManager = &RoomManager{
		rooms: make(map[string]*Room),
	}
	return roomManager
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
func (roomManager *RoomManager) HandleNewWebSocket(conn *WebSocketConn, request *http.Request) {
	util.Infof("On Open %v", request)
	defer func() {

		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	conn.On("message", func(message []byte) {
		defer func() {

			if r := recover(); r != nil {
				fmt.Printf("捕获到的错误：%s\n", r)
			}
		}()
		//解析Json数据
		request, err := util.Unmarshal(string(message))
		if err != nil {
			util.Errorf("解析Json数据Unmarshal错误 %v", err)
			return
		}
		//定义数据
		var data map[string]interface{} = nil
		tmp, found := request["data"]
		if !found {
			util.Errorf("没有发现数据!")
			return
		}
		data = tmp.(map[string]interface{})
		roomId := data["roomId"].(string)
		util.Infof("房间Id: %v", roomId)
		//根据roomId获取房间
		room := roomManager.getRoom(roomId)
		if room == nil {
			room = roomManager.createRoom(roomId)
		}

		switch request["type"] {
		case messagetype.JoinRoom.String():
			onJoinRoom(conn, data, room, roomManager)
			break
		case messagetype.Offer.String():
			fallthrough
		case messagetype.Answer.String():
			fallthrough
		case messagetype.Candidate.String():
			onCandidate(conn, data, room, roomManager, request)
			break
		case messagetype.HangUp.String():
			onHangUp(conn, data, room, roomManager, request)
			break
		default:
			{
				util.Warnf("未知的请求 %v", request)
			}
			break
		}

	})
	//连接关闭事件处理
	conn.On("close", func(code int, text string) {
		onClose(conn, roomManager)
	})

}

func onHangUp(conn *WebSocketConn, data map[string]interface{}, room *Room, manager *RoomManager, request map[string]interface{}) {
	sessionID := data["sessionId"].(string)
	ids := strings.Split(sessionID, "-")
	if user, ok := room.users[ids[0]]; !ok {
		util.Warnf("用户 [" + ids[0] + "] 没有找到")
		return
	} else {
		hangUp := map[string]interface{}{
			"type": messagetype.HangUp.String(),
			"data": map[string]interface{}{
				"to":        ids[0],
				"sessionId": sessionID,
			},
		}
		//发送信息给目标User,即自己[0]
		user.conn.Send(util.Marshal(hangUp))
	}

	if user, ok := room.users[ids[1]]; !ok {
		util.Warnf("用户 [" + ids[1] + "] 没有找到")
		return
	} else {
		hangUp := map[string]interface{}{
			"type": messagetype.HangUp.String(),
			"data": map[string]interface{}{
				"to":        ids[1],
				"sessionId": sessionID,
			},
		}
		//发送信息给目标User,即对方[1]
		user.conn.Send(util.Marshal(hangUp))
	}

}

//offer/answer/candidate消息处理
func onCandidate(conn *WebSocketConn, data map[string]interface{}, room *Room, manager *RoomManager, request map[string]interface{}) {
	to := data["to"].(string)
	if user, ok := room.users[to]; !ok {
		util.Errorf("没有发现用户[" + to + "]")
		return
	} else {
		user.conn.Send(util.Marshal(request))
	}

}

func onJoinRoom(conn *WebSocketConn, data map[string]interface{}, room *Room, roomManager *RoomManager) {
	//创建一个User
	user := User{
		conn: conn,
		info: UserInfo{
			ID:   data["id"].(string),
			Name: data["name"].(string),
		},
	}
	room.users[user.info.ID] = user
	roomManager.notifyUserUpdate(conn, room.users)
}

//通知所有的用户更新
func (roomManager *RoomManager) notifyUserUpdate(conn *WebSocketConn, users map[string]User) {
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
func onClose(conn *WebSocketConn, roomManager *RoomManager) {
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
