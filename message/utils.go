package message

type MessageType int32

const (
	MESSAGE_TYPE_ERROR MessageType = -1

	MESSAGE_TYPE_ACK   MessageType = 0
	MESSAGE_TYPE_USER  MessageType = 1
	MESSAGE_TYPE_GROUP MessageType = 2
)

func CheckMessageType(msg *Message) MessageType {
	switch msg.Type {
	case Type_Acknowledge:
		return MESSAGE_TYPE_ACK
	case Type_FRIEND_TEXT, Type_FRIEND_IMAGE, Type_FRIEND_FILE, Type_FRIEND_REQUEST, Type_FRIEND_ACCEPT, Type_FRIEND_DISBAND:
		return MESSAGE_TYPE_USER
	case Type_GROUP_TEXT, Type_GROUP_IMAGE, Type_GROUP_FILE, Type_GROUP_USER_CHANGE, Type_GROUP_REQUEST, Type_GROUP_ACCEPT, Type_GROUP_DISBAND:
		return MESSAGE_TYPE_GROUP
	default:
		return MESSAGE_TYPE_ERROR
	}
}
