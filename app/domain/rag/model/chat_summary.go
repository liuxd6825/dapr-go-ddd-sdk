package model

import (
	"time"
)

type ChatSummary struct {
	Base `bson:",inline"`

	ChatID             string `gorm:"chat_id;uniqueIndex" json:"chatId"`
	Summary            string `gorm:"summary" json:"summary"`
	CompressedPosition int    `gorm:"compressed_position" json:"compressedPosition"`
	MessageCount       int    `gorm:"message_count" json:"messageCount"`
	OriginalTokens     int    `gorm:"original_tokens" json:"originalTokens"`
	CompressedTokens   int    `gorm:"compressed_tokens" json:"compressedTokens"`
	IsLocked           bool   `gorm:"is_locked" json:"isLocked"`
}

func (c *ChatSummary) GetChatID() string {
	return c.ChatID
}

func (c *ChatSummary) SetChatID(chatId string) {
	c.ChatID = chatId
}

func (c *ChatSummary) GetSummary() string {
	return c.Summary
}

func (c *ChatSummary) SetSummary(summary string) {
	c.Summary = summary
}

func (c *ChatSummary) GetCompressedPosition() int {
	return c.CompressedPosition
}

func (c *ChatSummary) SetCompressedPosition(pos int) {
	c.CompressedPosition = pos
}

func (c *ChatSummary) GetIsLocked() bool {
	return c.IsLocked
}

func (c *ChatSummary) SetIsLocked(locked bool) {
	c.IsLocked = locked
}

func (c *ChatSummary) CanInsert() bool {
	return !c.IsLocked
}

func (c *ChatSummary) Lock() {
	c.IsLocked = true
}

func (c *ChatSummary) Unlock() {
	c.IsLocked = false
}

func (c *ChatSummary) GetCreatedTime() time.Time {
	if c.CreatedTime == nil {
		return time.Time{}
	}
	return *c.CreatedTime
}

func (c *ChatSummary) GetUpdatedTime() time.Time {
	if c.UpdatedTime == nil {
		return time.Time{}
	}
	return *c.UpdatedTime
}