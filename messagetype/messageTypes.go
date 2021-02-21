package messagetype

type MessageType int32

const (
	HeartPackage   MessageType = iota
	Login                      //登录
	Close                      //退出
	CheckOut                   //踢出
	JoinRoom                   //加入房间
	Offer                      //Offer消息
	Answer                     //Answer消息
	Candidate                  //Candidate消息
	HangUp                     //挂断
	LeaveRoom                  //离开房间
	UpdateUserList             //更新房间用户列表

)

func (p MessageType) String() string {
	switch p {
	case Close:
		return "close"
	case Login:
		return "login"
	case CheckOut:
		return "checkOut"
	case HeartPackage:
		return "heartPackage"
	case JoinRoom:
		return "joinRoom"
	case Offer:
		return "offer"
	case Answer:
		return "answer"
	case Candidate:
		return "candidate"
	case HangUp:
		return "hangUp"
	case LeaveRoom:
		return "leaveRoom"
	case UpdateUserList:
		return "updateUserList"
	default:
		return "UNKNOWN"
	}
}
