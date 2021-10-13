package requests

type Message struct {
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	PostType    string `json:"post_type"`    // message
	MessageType string `json:"message_type"` // group
	SubType     string `json:"sub_type"`     // normal、anonymous、notice	消息子类型, 正常消息是 normal, 匿名消息是 anonymous, 系统提示 ( 如「管理员已禁止群内匿名聊天」 ) 是 notice
	MessageId   int32  `json:"message_id"`
	GroupId     int64  `json:"group_id"`
	UserId      int64  `json:"user_id"`
	RawMessage  string `json:"raw_message"`
}
