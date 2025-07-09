package interfaces

import "github.com/baoerzuikeai/QQBot-NapCat/internal/repository/sqlite"

type AIHistoryRepository interface {
	SaveAIhistory(userID int64, content, role, sessionID string) error
	GetAIHistoryByUserID(userID int64) ([]sqlite.AIHistory, error)
}
