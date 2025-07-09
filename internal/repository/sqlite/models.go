package sqlite

type AIHistory struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`              // User's prompt or question
	Content   string `json:"ontent"`               // AI's response content
	Role      string `json:"ole"`                  // Role of the AI in the response (e.g., "assistant")
	Timestamp int64  `json:"timestamp"`            // Timestamp of the chat
	SessionID string `json:"session_id,omitempty"` // Session ID for grouping related messages
}
