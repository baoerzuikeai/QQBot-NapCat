package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type DeepSeekClient struct {
	apiKey string
	apiURL string
}

type ChatResponse struct {
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewDeepSeekClient(apiKey, apiURL string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey: apiKey,
		apiURL: apiURL,
	}
}

func (c *DeepSeekClient) GetResponse(ctx context.Context, history []Message, question string) ([]Choice, error) {
	// 构造请求体
	history = append(history, Message{
		Role:    "user",
		Content: question,
	})
	requestBody, err := json.Marshal(map[string]interface{}{
		"messages": history,
		"model":    "deepseek-chat",
	})
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("DeepSeek 服务返回错误: " + resp.Status)
	}
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Choices, nil
}
