package service

import (
	"context"

	"github.com/baoerzuikeai/QQBot-NapCat/internal/aiclient"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/repository/interfaces"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/repository/sqlite"
)

type AIChatServiceInterface interface {
	GetResponse(history []aiclient.Message, question string) ([]aiclient.Choice, error)
	SaveAIhistory(userID int64, content, role, sessionID string) error
	GetAIHistoryByUserID(userID int64) ([]sqlite.AIHistory, error)
}

type AIChatService struct {
	aiHistoryRepo interfaces.AIHistoryRepository
	aiClient      *aiclient.DeepSeekClient
}

func NewAIChatService(aiHistoryRepo interfaces.AIHistoryRepository, aiClient *aiclient.DeepSeekClient) *AIChatService {
	return &AIChatService{
		aiHistoryRepo: aiHistoryRepo,
		aiClient:      aiClient,
	}
}

func (s *AIChatService) GetResponse(history []aiclient.Message, question string) ([]aiclient.Choice, error) {
	return s.aiClient.GetResponse(context.Background(), history, question)
}

func (s *AIChatService) SaveAIhistory(userID int64, content, role, sessionID string) error {
	return s.aiHistoryRepo.SaveAIhistory(userID, content, role, sessionID)
}

func (s *AIChatService) GetAIHistoryByUserID(userID int64) ([]sqlite.AIHistory, error) {
	return s.aiHistoryRepo.GetAIHistoryByUserID(userID)
}
