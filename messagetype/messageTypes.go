package messagetype

type MessageType int32

const (
	HeartPackage MessageType = iota
)

func (p MessageType) String() string {
	switch p {
	case HeartPackage:
		return "heartPackage"
	default:
		return "UNKNOWN"
	}
}
