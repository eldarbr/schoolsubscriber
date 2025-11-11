package tgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eldarbr/schoolsubscriber/internal/myerrs"
)

type TGBot struct {
	apiKey string
	chatID int64
}

const baseUrl = "https://api.telegram.org/"

func (bot TGBot) SendMessage(ctx context.Context, msg string) error {
	const methodName = "sendMessage"

	fullUrl, err := url.JoinPath(baseUrl, "bot"+bot.apiKey, methodName)
	if err != nil {
		return fmt.Errorf("JoinPath: %w", err)
	}

	message := Message{
		Text:   msg,
		ChatID: bot.chatID,
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(jsonMsg))
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &myerrs.StatusCodeError{StatusCode: response.StatusCode}
	}

	return nil
}

func NewBot(key string, chatID int64) TGBot {
	return TGBot{
		apiKey: key,
		chatID: chatID,
	}
}
