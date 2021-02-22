package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"p2pserve/messagetype"
	"p2pserve/util"
)

type ALLManager struct {
	users    map[string]*User
	sessions map[string]*Session
}

func (r *ALLManager) JoinUser(userId string, roomId string, conn *WebSocketConn) {
	util.Infof("login" + userId)
	curUser, ok := r.users[userId]
	if ok {
		curUser.conn.Send(util.GenMessage(messagetype.CheckOut.String(), gin.H{
			"addr": conn.socket.RemoteAddr().String(),
		}))
		curUser.conn.Emit("close", -1, "check_out")
	}
	r.users[userId] = &User{
		conn: conn,
		info: UserInfo{
			ID:     userId,
			Name:   userId,
			RoomId: roomId,
		},
	}
	r.notifyUserList()

}

func (r *ALLManager) HandleNewWebSocket(conn *WebSocketConn, request *http.Request) {

	conn.On("message", func(message []byte) {
		request, err := util.Unmarshal(string(message))
		if err != nil {
			util.Errorf("解析Json数据Unmarshal错误 %v", err)
			return
		}
		switch request["type"] {
		case messagetype.JoinRoom.String():
			toUserId := request["to"].(string)
			toUser, ok := r.users[toUserId]
			if ok {
				toUser.conn.Send(util.GenMessage(messagetype.JoinRoom.String(), request))
			}
			break
		default:
			break

		}

		fmt.Println(request)
	})
	conn.On("close", func(code int, text string) {
		userId := ""
		for _, u := range r.users {
			if u.conn == conn {
				userId = u.info.ID
			}
		}
		if userId == "" {
			util.Errorf("没有查找到用户")
			return
		}
		delete(r.users, userId)
		r.notifyUserUpdate(util.GenMessage(messagetype.Close.String(), gin.H{"userId": userId, "text": text, "code": code}))
		r.notifyUserList()
	})
	userId := request.URL.Query().Get("userId")
	roomId := request.URL.Query().Get("roomId")
	util.Infof(roomId)
	r.JoinUser(userId, roomId, conn)

}

func (r *ALLManager) notifyUserUpdate(message string) {
	for _, user := range r.users {
		user.conn.Send(message)
	}
}

func (r *ALLManager) notifyUserList() {
	infos := []UserInfo{}
	for _, userClient := range r.users {
		infos = append(infos, userClient.info)
	}
	request := make(map[string]interface{})
	request["type"] = messagetype.UpdateAllUserList.String()
	//数据
	request["data"] = infos
	//迭代所有的User
	for _, user := range r.users {
		//将Json数据发送给每一个User
		user.conn.Send(util.Marshal(request))
	}
}

func NewAllManager() *ALLManager {
	aLLManager := new(ALLManager)
	aLLManager.users = make(map[string]*User)
	aLLManager.sessions = make(map[string]*Session)
	return aLLManager
}
