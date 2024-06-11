package notifier

import (
	"os"
	"testing"
)

func TestTelegramNotifier_Notify(t *testing.T) {
	// Set the token and chat ID for testing
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	// Call the Notify function
	err := TelegramNotifier.Notify(chatID, "Hello, World!", token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
