package ilert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// GetChatChannelThreadRepliesInput represents the input of a GetChatChannelThreadReplies operation.
type GetChatChannelThreadRepliesInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the replies belong to
	ThreadID *int64

	// an integer specifying the starting point (beginning with 0) when paging through a list of entities
	StartIndex *int

	// the maximum number of results when paging through a list of entities.
	// Maximum: 100
	MaxResults *int
}

// GetChatChannelThreadRepliesOutput represents the output of a GetChatChannelThreadReplies operation.
type GetChatChannelThreadRepliesOutput struct {
	_                        struct{}
	ChatChannelThreadReplies []*ChatChannelThreadReply
}

// GetChatChannelThreadReplies lists the replies of a thread, newest first.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) GetChatChannelThreadReplies(input *GetChatChannelThreadRepliesInput) (*GetChatChannelThreadRepliesOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)
	if input.StartIndex != nil {
		q.Add("start-index", strconv.Itoa(*input.StartIndex))
	}
	if input.MaxResults != nil {
		q.Add("max-results", strconv.Itoa(*input.MaxResults))
	}

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/messages/%d/thread-replies?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	replies := make([]*ChatChannelThreadReply, 0)
	err = json.Unmarshal(resp.Body(), &replies)
	if err != nil {
		return nil, err
	}

	return &GetChatChannelThreadRepliesOutput{ChatChannelThreadReplies: replies}, nil
}

// GetChatChannelThreadReplyInput represents the input of a GetChatChannelThreadReply operation.
type GetChatChannelThreadReplyInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the reply belongs to
	ThreadID *int64

	ChatChannelThreadReplyID *int64
}

// GetChatChannelThreadReplyOutput represents the output of a GetChatChannelThreadReply operation.
type GetChatChannelThreadReplyOutput struct {
	_                      struct{}
	ChatChannelThreadReply *ChatChannelThreadReply
}

// GetChatChannelThreadReply gets the thread reply with the specified id.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) GetChatChannelThreadReply(input *GetChatChannelThreadReplyInput) (*GetChatChannelThreadReplyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReplyID == nil {
		return nil, errors.New("chat channel thread reply id is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/%d/messages/%d/thread-replies/%d?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, *input.ChatChannelThreadReplyID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &GetChatChannelThreadReplyOutput{ChatChannelThreadReply: reply}, nil
}

// CreateChatChannelThreadReplyInput represents the input of a CreateChatChannelThreadReply operation.
type CreateChatChannelThreadReplyInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message being replied to. The first reply turns that message into a thread.
	ThreadID *int64

	// the reply to write. The author is the user the API key or access token belongs to.
	ChatChannelThreadReply *ChatChannelThreadReply
}

// CreateChatChannelThreadReplyOutput represents the output of a CreateChatChannelThreadReply operation.
type CreateChatChannelThreadReplyOutput struct {
	_                      struct{}
	ChatChannelThreadReply *ChatChannelThreadReply
}

// CreateChatChannelThreadReply replies to a message. The first reply turns that message into a
// thread, and the author is the user the API key or access token belongs to.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) CreateChatChannelThreadReply(input *CreateChatChannelThreadReplyInput) (*CreateChatChannelThreadReplyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReply == nil {
		return nil, errors.New("chat channel thread reply input is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(input.ChatChannelThreadReply).
		Post(fmt.Sprintf("%s/%d/messages/%d/thread-replies?%s",
			apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 201); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &CreateChatChannelThreadReplyOutput{ChatChannelThreadReply: reply}, nil
}

// UpdateChatChannelThreadReplyInput represents the input of an UpdateChatChannelThreadReply operation.
type UpdateChatChannelThreadReplyInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the reply belongs to
	ThreadID *int64

	ChatChannelThreadReplyID *int64
	ChatChannelThreadReply   *ChatChannelThreadReply
}

// UpdateChatChannelThreadReplyOutput represents the output of an UpdateChatChannelThreadReply operation.
type UpdateChatChannelThreadReplyOutput struct {
	_                      struct{}
	ChatChannelThreadReply *ChatChannelThreadReply
}

// UpdateChatChannelThreadReply updates an existing thread reply. Only its author may update it.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) UpdateChatChannelThreadReply(input *UpdateChatChannelThreadReplyInput) (*UpdateChatChannelThreadReplyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReplyID == nil {
		return nil, errors.New("chat channel thread reply id is required")
	}
	if input.ChatChannelThreadReply == nil {
		return nil, errors.New("chat channel thread reply input is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(input.ChatChannelThreadReply).
		Put(fmt.Sprintf("%s/%d/messages/%d/thread-replies/%d?%s",
			apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, *input.ChatChannelThreadReplyID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &UpdateChatChannelThreadReplyOutput{ChatChannelThreadReply: reply}, nil
}

// DeleteChatChannelThreadReplyInput represents the input of a DeleteChatChannelThreadReply operation.
type DeleteChatChannelThreadReplyInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the reply belongs to
	ThreadID *int64

	ChatChannelThreadReplyID *int64
}

// DeleteChatChannelThreadReplyOutput represents the output of a DeleteChatChannelThreadReply operation.
type DeleteChatChannelThreadReplyOutput struct {
	_ struct{}

	// the deleted reply, as it remains in the list
	ChatChannelThreadReply *ChatChannelThreadReply
}

// DeleteChatChannelThreadReply deletes a thread reply. Only its author may delete it, and the
// delete is soft, so this operation answers 200 with the reply rather than 204. Deleting a reply
// that is already deleted is rejected rather than treated as a no-op.
// https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) DeleteChatChannelThreadReply(input *DeleteChatChannelThreadReplyInput) (*DeleteChatChannelThreadReplyOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReplyID == nil {
		return nil, errors.New("chat channel thread reply id is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d/messages/%d/thread-replies/%d?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, *input.ChatChannelThreadReplyID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &DeleteChatChannelThreadReplyOutput{ChatChannelThreadReply: reply}, nil
}

// AddChatChannelThreadReplyReactionInput represents the input of an AddChatChannelThreadReplyReaction operation.
type AddChatChannelThreadReplyReactionInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the reply belongs to
	ThreadID *int64

	ChatChannelThreadReplyID *int64

	// the reaction code to add on behalf of the authenticated user
	Code *string
}

// AddChatChannelThreadReplyReactionOutput represents the output of an AddChatChannelThreadReplyReaction operation.
type AddChatChannelThreadReplyReactionOutput struct {
	_                      struct{}
	ChatChannelThreadReply *ChatChannelThreadReply
}

// AddChatChannelThreadReplyReaction adds a reaction of the authenticated user to a thread reply and
// returns the updated reply. Reacting twice with the same code, or exceeding the limit of codes per
// reply, answers 400. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) AddChatChannelThreadReplyReaction(input *AddChatChannelThreadReplyReactionInput) (*AddChatChannelThreadReplyReactionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReplyID == nil {
		return nil, errors.New("chat channel thread reply id is required")
	}
	if input.Code == nil {
		return nil, errors.New("reaction code is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().SetBody(ChatChannelMessageReaction{Code: *input.Code}).
		Post(fmt.Sprintf("%s/%d/messages/%d/thread-replies/%d/reactions?%s",
			apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, *input.ChatChannelThreadReplyID, q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &AddChatChannelThreadReplyReactionOutput{ChatChannelThreadReply: reply}, nil
}

// RemoveChatChannelThreadReplyReactionInput represents the input of a RemoveChatChannelThreadReplyReaction operation.
type RemoveChatChannelThreadReplyReactionInput struct {
	_ struct{}

	// the id of the channel the thread belongs to, for the channel type "ALERT" this is
	// the id of the alert
	ChannelID *int64

	// the type of the channel, required
	// possible values: "ALERT"
	ChannelType *string

	// the id of the message the reply belongs to
	ThreadID *int64

	ChatChannelThreadReplyID *int64

	// the reaction code to remove on behalf of the authenticated user
	Code *string
}

// RemoveChatChannelThreadReplyReactionOutput represents the output of a RemoveChatChannelThreadReplyReaction operation.
type RemoveChatChannelThreadReplyReactionOutput struct {
	_                      struct{}
	ChatChannelThreadReply *ChatChannelThreadReply
}

// RemoveChatChannelThreadReplyReaction removes the reaction of the authenticated user from a thread
// reply and returns the updated reply. https://docs.ilert.com/developer-docs/rest-api/api-reference/chat-messages
func (c *Client) RemoveChatChannelThreadReplyReaction(input *RemoveChatChannelThreadReplyReactionInput) (*RemoveChatChannelThreadReplyReactionOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.ChannelID == nil {
		return nil, errors.New("channel id is required")
	}
	if input.ChannelType == nil {
		return nil, errors.New("channel type is required")
	}
	if input.ThreadID == nil {
		return nil, errors.New("thread id is required")
	}
	if input.ChatChannelThreadReplyID == nil {
		return nil, errors.New("chat channel thread reply id is required")
	}
	if input.Code == nil {
		return nil, errors.New("reaction code is required")
	}

	q := url.Values{}
	q.Add("channel-type", *input.ChannelType)

	resp, err := c.httpClient.R().Delete(fmt.Sprintf("%s/%d/messages/%d/thread-replies/%d/reactions/%s?%s",
		apiRoutes.chatChannels, *input.ChannelID, *input.ThreadID, *input.ChatChannelThreadReplyID,
		url.PathEscape(*input.Code), q.Encode()))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	reply := &ChatChannelThreadReply{}
	err = json.Unmarshal(resp.Body(), reply)
	if err != nil {
		return nil, err
	}

	return &RemoveChatChannelThreadReplyReactionOutput{ChatChannelThreadReply: reply}, nil
}
