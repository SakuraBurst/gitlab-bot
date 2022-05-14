package models

import "fmt"

type TelegramResponse struct {
	Ok     bool           `json:"ok"`
	Result TelegramResult `json:"result"`
}

type TelegramResult struct {
	MessageID int                `json:"message_id"`
	From      TelegramResultFrom `json:"from"`
	Chat      TelegramResultChat `json:"chat"`
	Date      int                `json:"date"`
	Text      string             `json:"text"`
}
type TelegramResultFrom struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type TelegramResultChat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
	Type      string `json:"type"`
}

type TelegramError struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (t TelegramError) Error() string {
	return fmt.Sprintf("Code %d, Message: %s", t.ErrorCode, t.Description)
}
