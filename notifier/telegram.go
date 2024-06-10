package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type telegramNotifier struct {
	client *http.Client
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var TelegramNotifier *telegramNotifier = &telegramNotifier{
	client: &http.Client{
		Timeout: 30 * time.Second,
	},
}

func (t *telegramNotifier) Notify(chatID, text, token string) error {
	message := telegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.telegram.org",
			Path:   fmt.Sprintf("/bot%s/sendMessage", token),
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(messageBytes)),
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body) // Read the response body for additional error info
		return fmt.Errorf("failed to send message to telegram: %s - %s", resp.Status, body)
	}

	return nil
}
