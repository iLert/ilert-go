package ilert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// channel-type is required by the API, which answers 400 without it, so every chat request
// has to carry it no matter which operation built the URL. A reaction is sent as its code
// alone, the counts and the reacting users are filled in by the API.
func TestChatChannelRequestsAlwaysSendChannelType(t *testing.T) {
	c := func(url string) *Client { return newTestClient(t, url) }
	message := &ChatChannelMessage{Content: "hi", ContentType: ChatChannelMessageContentType.Text}
	cases := []struct {
		name   string
		call   func(*Client) error
		method string
		want   string
		body   string
	}{
		{
			name: "list messages",
			call: func(cl *Client) error {
				_, err := cl.GetChatChannelMessages(&GetChatChannelMessagesInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert)})
				return err
			},
			method: http.MethodGet,
			want:   "/api/chat-channels/7/messages",
		},
		{
			name: "get message",
			call: func(cl *Client) error {
				_, err := cl.GetChatChannelMessage(&GetChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9)})
				return err
			},
			method: http.MethodGet,
			want:   "/api/chat-channels/7/messages/9",
		},
		{
			name: "create message",
			call: func(cl *Client) error {
				_, err := cl.CreateChatChannelMessage(&CreateChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessage: message})
				return err
			},
			method: http.MethodPost,
			want:   "/api/chat-channels/7/messages",
			body:   `{"content":"hi","contentType":"TEXT"}`,
		},
		{
			name: "update message",
			call: func(cl *Client) error {
				_, err := cl.UpdateChatChannelMessage(&UpdateChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9), ChatChannelMessage: message})
				return err
			},
			method: http.MethodPut,
			want:   "/api/chat-channels/7/messages/9",
			body:   `{"content":"hi","contentType":"TEXT"}`,
		},
		{
			name: "delete message",
			call: func(cl *Client) error {
				_, err := cl.DeleteChatChannelMessage(&DeleteChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9)})
				return err
			},
			method: http.MethodDelete,
			want:   "/api/chat-channels/7/messages/9",
		},
		{
			name: "add reaction",
			call: func(cl *Client) error {
				_, err := cl.AddChatChannelMessageReaction(&AddChatChannelMessageReactionInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9), Code: String("thumbsup")})
				return err
			},
			method: http.MethodPost,
			want:   "/api/chat-channels/7/messages/9/reactions",
			body:   `{"code":"thumbsup"}`,
		},
		{
			name: "remove reaction",
			call: func(cl *Client) error {
				_, err := cl.RemoveChatChannelMessageReaction(&RemoveChatChannelMessageReactionInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9), Code: String("thumbsup")})
				return err
			},
			method: http.MethodDelete,
			want:   "/api/chat-channels/7/messages/9/reactions/thumbsup",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var method, path, body string
			var query url.Values
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				path = r.URL.Path
				query = r.URL.Query()
				raw, _ := io.ReadAll(r.Body)
				body = string(raw)
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet && r.URL.Path == "/api/chat-channels/7/messages" {
					_, _ = w.Write([]byte(`[]`))
					return
				}
				if r.Method == http.MethodPost && r.URL.Path == "/api/chat-channels/7/messages" {
					w.WriteHeader(http.StatusCreated)
				}
				_, _ = w.Write([]byte(`{"id":9,"content":"hi","contentType":"TEXT"}`))
			}))
			defer srv.Close()

			if err := tc.call(c(srv.URL)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if method != tc.method {
				t.Errorf("method = %s, want %s", method, tc.method)
			}
			if path != tc.want {
				t.Errorf("path = %q, want %q", path, tc.want)
			}
			if got := query.Get("channel-type"); got != "ALERT" {
				t.Errorf("channel-type = %q, want ALERT", got)
			}
			if tc.body != "" && body != tc.body {
				t.Errorf("body = %s, want %s", body, tc.body)
			}
		})
	}
}

// A reaction code is user input in the path, so it has to be escaped into a single segment.
func TestRemoveChatChannelMessageReactionEscapesCode(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":9,"content":"hi","contentType":"TEXT"}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).RemoveChatChannelMessageReaction(&RemoveChatChannelMessageReactionInput{
		ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9),
		Code: String("thumbs up/down"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/chat-channels/7/messages/9/reactions/thumbs%20up%2Fdown" {
		t.Errorf("path = %q, want the code escaped into a single segment", path)
	}
}

// The delete is soft, so unlike every other delete in this SDK it answers 200 with the
// message body rather than 204, and the caller gets the tombstone back.
func TestDeleteChatChannelMessageReturnsTheTombstone(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":9,"content":"<Deleted on 2026-09-17T11:04:28Z>","contentType":"TEXT","deleted":true,"thread":true,"threadMessageCount":1}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).DeleteChatChannelMessage(&DeleteChatChannelMessageInput{
		ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ChatChannelMessageID: Int64(9),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.ChatChannelMessage.Deleted {
		t.Error("deleted = false, want true")
	}
	if result.ChatChannelMessage.ThreadMessageCount != 1 {
		t.Errorf("threadMessageCount = %d, want 1", result.ChatChannelMessage.ThreadMessageCount)
	}
}

// Reactions decode into their own type rather than a bare list of codes, since one entry
// carries the code plus everyone who used it.
func TestChatChannelMessageReactionsDecode(t *testing.T) {
	message := &ChatChannelMessage{}
	err := json.Unmarshal([]byte(`{"id":9,"reactions":[{"code":"thumbsup","count":2,"users":[1,2],"bots":[3]}]}`), message)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(message.Reactions) != 1 {
		t.Fatalf("reactions = %+v, want one entry", message.Reactions)
	}
	r := message.Reactions[0]
	if r.Code != "thumbsup" || r.Count != 2 || len(r.Users) != 2 || len(r.Bots) != 1 {
		t.Errorf("reaction = %+v, want thumbsup with 2 users and 1 bot", r)
	}
}

func TestChatChannelMessageRequiredInputs(t *testing.T) {
	c := NewClient()
	cases := map[string]func() error{
		"nil input": func() error { _, err := c.GetChatChannelMessages(nil); return err },
		"no channel id": func() error {
			_, err := c.GetChatChannelMessages(&GetChatChannelMessagesInput{ChannelType: String("ALERT")})
			return err
		},
		"no channel type": func() error {
			_, err := c.GetChatChannelMessages(&GetChatChannelMessagesInput{ChannelID: Int64(7)})
			return err
		},
		"no message id": func() error {
			_, err := c.GetChatChannelMessage(&GetChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String("ALERT")})
			return err
		},
		"no message body": func() error {
			_, err := c.CreateChatChannelMessage(&CreateChatChannelMessageInput{ChannelID: Int64(7), ChannelType: String("ALERT")})
			return err
		},
		"no reaction code": func() error {
			_, err := c.AddChatChannelMessageReaction(&AddChatChannelMessageReactionInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ChatChannelMessageID: Int64(9)})
			return err
		},
		"no code on remove": func() error {
			_, err := c.RemoveChatChannelMessageReaction(&RemoveChatChannelMessageReactionInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ChatChannelMessageID: Int64(9)})
			return err
		},
	}

	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Error("error = nil, want an error")
			}
		})
	}
}
