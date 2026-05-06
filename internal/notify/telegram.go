package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBase = "https://api.telegram.org"

// TelegramNotifier sends job alerts via the Telegram Bot API.
type TelegramNotifier struct {
	token  string
	chatID string
	apiBase string
}

// NewTelegramNotifier creates a TelegramNotifier with the given bot token and chat ID.
func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:   token,
		chatID:  chatID,
		apiBase: telegramAPIBase,
	}
}

type telegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// Notify sends a Telegram message describing the job result.
func (t *TelegramNotifier) Notify(jobName string, err error) error {
	text := fmt.Sprintf("✅ cronwrap: job %q succeeded.", jobName)
	if err != nil {
		text = fmt.Sprintf("❌ cronwrap: job %q failed: %v", jobName, err)
	}

	payload := telegramPayload{
		ChatID: t.chatID,
		Text:   text,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("telegram: marshal payload: %w", marshalErr)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBase, t.token)
	resp, httpErr := http.Post(url, "application/json", bytes.NewReader(body)) //nolint:noctx
	if httpErr != nil {
		return fmt.Errorf("telegram: send message: %w", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}

	return nil
}
