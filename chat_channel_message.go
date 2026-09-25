package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// ChatChannelMessage definition https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
type ChatChannelMessage struct {
	ID int64 `json:"id,omitempty"`

	// the message body, required on create and update
	Content string `json:"content"`

	// the type of the content, required on create and update
	// possible values: "TEXT"
	ContentType string `json:"contentType"`

	// who wrote the message, read-only
	// possible values: "USER", "BOT", "EXTERNAL"
	AuthorType string `json:"authorType,omitempty"`

	// the id of the author, read-only. Only set when the author is a user or a bot.
	AuthorID int64 `json:"authorId,omitempty"`

	// the name of the bot that wrote the message, read-only
	BotName string `json:"botName,omitempty"`

	// the name of the external author, read-only. Set when the message was submitted
	// with an event rather than written by a user or a bot.
	ExternalCreatorName string `json:"externalCreatorName,omitempty"`

	// how the message is meant to be rendered, read-only. Messages written by users always
	// carry "DEFAULT"; ilert uses further styles internally, for example for the intermediate
	// steps of its AI agents, and the API documents that only "DEFAULT" should be rendered.
	// possible values: "DEFAULT"
	MessageStyle string `json:"messageStyle,omitempty"`

	// whether the message was edited after it was written, read-only
	Edited bool `json:"edited,omitempty"`

	// whether the message was deleted, read-only. Deletes are soft: the message stays in the
	// list with this flag set and its content replaced by a placeholder.
	Deleted bool `json:"deleted,omitempty"`

	// whether the message resolved the alert of the channel, read-only
	ResolveComment bool `json:"resolveComment,omitempty"`

	// the reactions on this message, one entry per code, read-only.
	// Written through AddChatChannelMessageReaction and RemoveChatChannelMessageReaction.
	Reactions []ChatChannelMessageReaction `json:"reactions,omitempty"`

	// whether the message has replies, read-only. A message becomes a thread with its first reply.
	Thread bool `json:"thread,omitempty"`

	// the number of replies in the thread, read-only
	ThreadMessageCount int64 `json:"threadMessageCount,omitempty"`

	// when the most recent reply was written, read-only. Date time string in ISO format.
	LastReplyDate string `json:"lastReplyDate,omitempty"`

	// who wrote the most recent reply, read-only
	// possible values: "USER", "BOT", "EXTERNAL"
	LastReplyAuthorType string `json:"lastReplyAuthorType,omitempty"`

	// the id of the author of the most recent reply, read-only
	LastReplyAuthorID int64 `json:"lastReplyAuthorId,omitempty"`

	// the name of the bot that wrote the most recent reply, read-only
	LastReplyBotName string `json:"lastReplyBotName,omitempty"`

	Created string `json:"created,omitempty"` // Date time string in ISO format
	Updated string `json:"updated,omitempty"` // Date time string in ISO format
}

// ChatChannelThreadReply defines a reply to a chat channel message
type ChatChannelThreadReply struct {
	ID int64 `json:"id,omitempty"`

	// the reply body, required on create and update
	Content string `json:"content"`

	// the type of the content, required on create and update
	// possible values: "TEXT"
	ContentType string `json:"contentType"`

	// the id of the message this reply belongs to, read-only
	ThreadID int64 `json:"threadId,omitempty"`

	// who wrote the reply, read-only
	// possible values: "USER", "BOT", "EXTERNAL"
	AuthorType string `json:"authorType,omitempty"`

	// the id of the author, read-only. Only set when the author is a user or a bot.
	AuthorID int64 `json:"authorId,omitempty"`

	// the name of the bot that wrote the reply, read-only
	BotName string `json:"botName,omitempty"`

	// how the reply is meant to be rendered, read-only
	// possible values: "DEFAULT"
	MessageStyle string `json:"messageStyle,omitempty"`

	// whether the reply was edited after it was written, read-only
	Edited bool `json:"edited,omitempty"`

	// whether the reply was deleted, read-only. Deletes are soft, as on the message.
	Deleted bool `json:"deleted,omitempty"`

	// the reactions on this reply, one entry per code, read-only
	Reactions []ChatChannelMessageReaction `json:"reactions,omitempty"`

	Created string `json:"created,omitempty"` // Date time string in ISO format
	Updated string `json:"updated,omitempty"` // Date time string in ISO format
}

// ChatChannelMessageReaction defines the reactions of one code on a message or thread reply
type ChatChannelMessageReaction struct {
	// the reaction code, for example an emoji shortcode
	Code string `json:"code"`

	// how many users and bots reacted with this code
	Count int64 `json:"count,omitempty"`

	// the ids of the users who reacted
	Users []int64 `json:"users,omitempty"`

	// the ids of the bots that reacted
	Bots []int64 `json:"bots,omitempty"`
}

// ChatChannelType defines the channel types that carry messages
var ChatChannelType = struct {
	Alert string
}{
	Alert: "ALERT",
}

// ChatChannelTypeAll defines the chat channel type list
var ChatChannelTypeAll = []string{
	ChatChannelType.Alert,
}

// ChatChannelMessageAuthorType defines who wrote a message or a thread reply
var ChatChannelMessageAuthorType = struct {
	User     string
	Bot      string
	External string
}{
	User:     "USER",
	Bot:      "BOT",
	External: "EXTERNAL",
}

// ChatChannelMessageAuthorTypeAll defines the chat channel message author type list
var ChatChannelMessageAuthorTypeAll = []string{
	ChatChannelMessageAuthorType.User,
	ChatChannelMessageAuthorType.Bot,
	ChatChannelMessageAuthorType.External,
}

// ChatChannelMessageContentType defines the content types of a message or a thread reply
var ChatChannelMessageContentType = struct {
	Text string
}{
	Text: "TEXT",
}

// ChatChannelMessageContentTypeAll defines the chat channel message content type list
var ChatChannelMessageContentTypeAll = []string{
	ChatChannelMessageContentType.Text,
}

// ChatChannelMessageStyle defines how a message or a thread reply is meant to be rendered
var ChatChannelMessageStyle = struct {
	Default string
}{
	Default: "DEFAULT",
}

// ChatChannelMessageStyleAll defines the chat channel message style list
var ChatChannelMessageStyleAll = []string{
	ChatChannelMessageStyle.Default,
}

// GetChatChannelMessagesInput represents the input of a GetChatChannelMessages operation.
type GetChatChannelMessagesInput struct {
	_ struct{}

	// the id of the channel the messages belong to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetChatChannelMessagesOutput represents the output of a GetChatChannelMessages operation.
type GetChatChannelMessagesOutput struct {
	_                   struct{}
	ChatChannelMessages []*ChatChannelMessage
}

// GetChatChannelMessages lists the messages of a channel. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) GetChatChannelMessages(input *GetChatChannelMessagesInput) (*GetChatChannelMessagesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/messages?%s", apiRoutes.chatChannels, *input.ChannelID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	messages := make([]*ChatChannelMessage, 0)
	err = json.Unmarshal(resp.Body(), &messages)
	if err != nil {
		return nil, err
	}

	return &GetChatChannelMessagesOutput{ChatChannelMessages: messages}, nil
}

// GetChatChannelMessageInput represents the input of a GetChatChannelMessage operation.
type GetChatChannelMessageInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	ChatChannelMessageID *int64
}

// GetChatChannelMessageOutput represents the output of a GetChatChannelMessage operation.
type GetChatChannelMessageOutput struct {
	_                  struct{}
	ChatChannelMessage *ChatChannelMessage
}

// GetChatChannelMessage gets the message with the specified id. An unknown id is answered
// with a 400 rather than a 404, so it surfaces as a *BadRequestAPIError (verified against
// the API on 17.09.2026). https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) GetChatChannelMessage(input *GetChatChannelMessageInput) (*GetChatChannelMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessageID == nil {
		return nil, errors.New("chat channel message id is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/messages/%d?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ChatChannelMessageID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &GetChatChannelMessageOutput{ChatChannelMessage: message}, nil
}

// CreateChatChannelMessageInput represents the input of a CreateChatChannelMessage operation.
type CreateChatChannelMessageInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the message to write. The author is the user the API key or access token belongs to.
	ChatChannelMessage *ChatChannelMessage
}

// CreateChatChannelMessageOutput represents the output of a CreateChatChannelMessage operation.
type CreateChatChannelMessageOutput struct {
	_                  struct{}
	ChatChannelMessage *ChatChannelMessage
}

// CreateChatChannelMessage writes a new message to a channel. The author is the user the API
// key or access token belongs to. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) CreateChatChannelMessage(input *CreateChatChannelMessageInput) (*CreateChatChannelMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessage == nil {
		return nil, errors.New("chat channel message input is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(input.ChatChannelMessage).
		Post(fmt.Sprintf("%s/%d/messages?%s", apiRoutes.chatChannels, *input.ChannelID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &CreateChatChannelMessageOutput{ChatChannelMessage: message}, nil
}

// UpdateChatChannelMessageInput represents the input of an UpdateChatChannelMessage operation.
type UpdateChatChannelMessageInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	ChatChannelMessageID *int64
	ChatChannelMessage   *ChatChannelMessage
}

// UpdateChatChannelMessageOutput represents the output of an UpdateChatChannelMessage operation.
type UpdateChatChannelMessageOutput struct {
	_                  struct{}
	ChatChannelMessage *ChatChannelMessage
}

// UpdateChatChannelMessage updates an existing message. Only its author may update it, the API
// answers 403 otherwise. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) UpdateChatChannelMessage(input *UpdateChatChannelMessageInput) (*UpdateChatChannelMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessageID == nil {
		return nil, errors.New("chat channel message id is required")
	}
	if input.ChatChannelMessage == nil {
		return nil, errors.New("chat channel message input is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(input.ChatChannelMessage).
		Put(fmt.Sprintf("%s/%d/messages/%d?%s", apiRoutes.chatChannels, *input.ChannelID, *input.ChatChannelMessageID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &UpdateChatChannelMessageOutput{ChatChannelMessage: message}, nil
}

// DeleteChatChannelMessageInput represents the input of a DeleteChatChannelMessage operation.
type DeleteChatChannelMessageInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	ChatChannelMessageID *int64
}

// DeleteChatChannelMessageOutput represents the output of a DeleteChatChannelMessage operation.
type DeleteChatChannelMessageOutput struct {
	_ struct{}

	// the deleted message, as it remains in the list
	ChatChannelMessage *ChatChannelMessage
}

// DeleteChatChannelMessage deletes a message. Only its author may delete it, and the delete is
// soft: the message stays in the list with Deleted set and its content replaced by a placeholder,
// which is why this operation answers 200 with the message rather than 204. Deleting a message
// that is already deleted is rejected rather than treated as a no-op.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) DeleteChatChannelMessage(input *DeleteChatChannelMessageInput) (*DeleteChatChannelMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessageID == nil {
		return nil, errors.New("chat channel message id is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d/messages/%d?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ChatChannelMessageID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &DeleteChatChannelMessageOutput{ChatChannelMessage: message}, nil
}

// AddChatChannelMessageReactionInput represents the input of an AddChatChannelMessageReaction operation.
type AddChatChannelMessageReactionInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	ChatChannelMessageID *int64

	// the reaction code to add on behalf of the authenticated user
	Code *string
}

// AddChatChannelMessageReactionOutput represents the output of an AddChatChannelMessageReaction operation.
type AddChatChannelMessageReactionOutput struct {
	_                  struct{}
	ChatChannelMessage *ChatChannelMessage
}

// AddChatChannelMessageReaction adds a reaction of the authenticated user to a message and returns
// the updated message. Reacting twice with the same code, or exceeding the limit of codes per
// message, answers 400. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) AddChatChannelMessageReaction(input *AddChatChannelMessageReactionInput) (*AddChatChannelMessageReactionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessageID == nil {
		return nil, errors.New("chat channel message id is required")
	}
	if input.Code == nil {
		return nil, errors.New("reaction code is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(ChatChannelMessageReaction{Code: *input.Code}).
		Post(fmt.Sprintf("%s/%d/messages/%d/reactions?%s",
			apiRoutes.chatChannels, *input.ChannelID, *input.ChatChannelMessageID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &AddChatChannelMessageReactionOutput{ChatChannelMessage: message}, nil
}

// RemoveChatChannelMessageReactionInput represents the input of a RemoveChatChannelMessageReaction operation.
type RemoveChatChannelMessageReactionInput struct {
	_ struct{}

	// the id of the channel the message belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	ChatChannelMessageID *int64

	// the reaction code to remove on behalf of the authenticated user
	Code *string
}

// RemoveChatChannelMessageReactionOutput represents the output of a RemoveChatChannelMessageReaction operation.
type RemoveChatChannelMessageReactionOutput struct {
	_                  struct{}
	ChatChannelMessage *ChatChannelMessage
}

// RemoveChatChannelMessageReaction removes the reaction of the authenticated user from a message
// and returns the updated message. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) RemoveChatChannelMessageReaction(input *RemoveChatChannelMessageReactionInput) (*RemoveChatChannelMessageReactionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ChatChannelMessageID == nil {
		return nil, errors.New("chat channel message id is required")
	}
	if input.Code == nil {
		return nil, errors.New("reaction code is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d/messages/%d/reactions/%s?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ChatChannelMessageID, url.PathEscape(*input.Code), q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	message := &ChatChannelMessage{}
	err = json.Unmarshal(resp.Body(), message)
	if err != nil {
		return nil, err
	}

	return &RemoveChatChannelMessageReactionOutput{ChatChannelMessage: message}, nil
}
