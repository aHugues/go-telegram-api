// package bot contains everything related to handling a bot
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ahugues/go-telegram-api/baseclt"
	"github.com/ahugues/go-telegram-api/servererror"
	"github.com/ahugues/go-telegram-api/structs"
)

type sentMessage struct {
	ChatID     int64  `json:"chat_id"`
	Text       string `json:"text"`
	Formatting string `json:"parse_mode,omitempty"`
}

type userResponse struct {
	OK   bool         `json:"ok"`
	User structs.User `json:"result"`
}

type Bot interface {
	// GetMe returns all the information concerning the current bot
	GetMe(ctx context.Context) (structs.User, error)

	// SendMessage sends the given message to the given chatID
	SendMessage(ctx context.Context, chatID int64, content string, formattingMode structs.FormattingOption) error
}

type ConcreteBot struct {
	token   string
	httpClt *http.Client
	baseURL string
}

func New(token string) *ConcreteBot {
	return &ConcreteBot{
		token:   token,
		httpClt: http.DefaultClient,
		baseURL: baseclt.BaseTelegramAPIURL,
	}
}

func (c *ConcreteBot) GetMe(ctx context.Context) (usr structs.User, err error) {
	url := fmt.Sprintf("%s/bot%s/getMe", c.baseURL, c.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return usr, fmt.Errorf("error building request: %w", err)
	}
	resp, err := c.httpClt.Do(req)
	if err != nil {
		return usr, fmt.Errorf("error sending request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return usr, servererror.FromResponse(resp)
	}

	var rawResponse userResponse
	if err := baseclt.ParseJSONBody(resp, &rawResponse); err != nil {
		return usr, fmt.Errorf("error parsing response: %w", err)
	}
	return rawResponse.User, nil
}

func (c *ConcreteBot) SendMessage(ctx context.Context, chatID int64, content string, formattingMode structs.FormattingOption) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", baseclt.BaseTelegramAPIURL, c.token)
	payload := sentMessage{
		ChatID:     chatID,
		Text:       content,
		Formatting: formattingMode.String(),
	}

	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error building payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bytesPayload))
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	req.Header.Add("content-type", "application/json")
	resp, err := c.httpClt.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return servererror.FromResponse(resp)
	}
	return nil
}
