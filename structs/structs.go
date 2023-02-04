// package structs contains all the basic types used by the Telegram API
package structs

import "reflect"

// User represents a User or a Bot
type User struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

// Chat represent a chat, either group or private
type Chat struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"` // Can be "private", "group", "supergroup" or "channel"
	Title       string `json:"title"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Bio         string `json:"bio"`
	Description string `json:"description"`
}

// Message represents a sent or received message
type Message struct {
	ID              int64  `json:"message_id"`
	From            User   `json:"from"`
	SenderChat      Chat   `json:"sender_chat"`
	Date            int64  `json:"date"` // Unix format
	Chat            Chat   `json:"chat"`
	ForwardFromChat Chat   `json:"forward_from_chat"`
	ViaBot          bool   `json:"via_bot"`
	Text            string `json:"text"`
}

// UpdateType represents the type of an update for future filters
type UpdateType int

const (
	UpdateMessage UpdateType = iota
	UpdateEditedMessage
	UpdateChannelPost
	UpdateEditedChannelPost
	UpdatedUnknown
)

func (e UpdateType) String() string {
	switch e {
	case UpdateMessage:
		return "message"
	case UpdateEditedMessage:
		return "edited_message"
	case UpdateChannelPost:
		return "channel_post"
	case UpdateEditedChannelPost:
		return "edited_channel_post"
	default:
		return "unknown"
	}
}

// Update represents an update gotten from the getUpdate method
type Update struct {
	ID                int64   `json:"update_id"`
	Message           Message `json:"message"`
	EditedMessage     Message `json:"edited_message"`
	ChannelPost       Message `json:"channel_post"`
	EditedChannelPost Message `json:"edited_channel_post"`
}

// Check every field in the update struct to check which one is not empty, and return it
//
// It is supposed that only a single type of update can be returned
func (u *Update) Type() UpdateType {
	var emptyMessage Message
	if !reflect.DeepEqual(u.Message, emptyMessage) {
		return UpdateMessage
	}
	if !reflect.DeepEqual(u.EditedMessage, emptyMessage) {
		return UpdateEditedMessage
	}
	if !reflect.DeepEqual(u.ChannelPost, emptyMessage) {
		return UpdateChannelPost
	}
	if !reflect.DeepEqual(u.EditedChannelPost, emptyMessage) {
		return UpdateEditedChannelPost
	}
	return UpdatedUnknown
}

type Updates struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type FormattingOption int

const (
	FormattingMarkdownV2 = iota
	FormattingHTML
	FormattingLegacy
	FormattingNone
)

func (o FormattingOption) String() string {
	switch o {
	case FormattingMarkdownV2:
		return "MarkdownV2"
	case FormattingHTML:
		return "HTML"
	case FormattingLegacy:
		return "Markdown"
	case FormattingNone:
		return ""
	default:
		return ""
	}
}
