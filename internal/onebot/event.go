package onebot

type BaseEvent struct {
	PostType string `json:"post_type"` // "message", "meta_event", "notice", "request"
	SelfID   int64  `json:"self_id"`   // The bot's self ID
	Time     int64  `json:"time"`      // Timestamp of the event
}

type Sender struct {
	UserID   int64  `json:"user_id"`  // ID of the user who
	NickName string `json:"nickname"` // Nickname of the user
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`   // Age of the user
	Card     string `json:"card"`  // Card name of the user (if applicable)
	Title    string `json:"title"` // Title of the user (if applicable, e.g., in a group)
	Level    string `json:"level"` // User's level (if applicable)
}

type MessageEvent struct {
	BaseEvent
	MessageType string  `json:"message_type"` // "private" or "group"
	MessageID   int64   `json:"message_id"`   // Unique ID for the message
	UserID      int64   `json:"user_id"`      // ID of the user who sent the message
	Message     Message `json:"message"`      // The message content
	RawMessage  string  `json:"raw_message"`  // Raw message content
	Font        int32   `json:"font"`         // Font ID (if applicable)
	Sender      Sender  `json:"sender"`       // Information about the sender
}

type PrivateMessageEvent struct {
	MessageEvent
	SubType string `json:"sub_type"` // "friend" or "group"
}

type GroupMessageEvent struct {
	MessageEvent
	GroupID   int64       `json:"group_id"`  // ID of the group where the message was sent
	Anonymous interface{} `json:"anonymous"` // Whether the message is anonymous
}
