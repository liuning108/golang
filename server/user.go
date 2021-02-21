package server

//用户信息
type UserInfo struct {
	ID   string `json:"id"`   //ID
	Name string `json:"name"` //名称
}

type User struct {
	//用户信息
	info UserInfo
	//连接对象
	conn *WebSocketConn
}

//会话信息
type Session struct {
	//会话id
	id string
	//消息来源
	from User
	//消息要发送的目标
	to User
}
