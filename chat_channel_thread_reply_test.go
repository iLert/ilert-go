package ilert

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// A thread reply lives under the message it answers, so every reply operation has to build
// that nested path and carry the channel-type the API requires, like the messages do.
func TestChatChannelThreadReplyRequestsAlwaysSendChannelType(t *testing.T) {
	c := func(url string) *Client { return newTestClient(t, url) }
	reply := &ChatChannelThreadReply{Content: "on it", ContentType: ChatChannelMessageContentType.Text}
	cases := []struct {
		name   string
		call   func(*Client) error
		method string
		want   string
		body   string
	}{
		{
			name: "list replies",
			call: func(cl *Client) error {
				_, err := cl.GetChatChannelThreadReplies(&GetChatChannelThreadRepliesInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9)})
				return err
			},
			method: http.MethodGet,
			want:   "/api/chat-channels/7/messages/9/thread-replies",
		},
		{
			name: "get reply",
			call: func(cl *Client) error {
				_, err := cl.GetChatChannelThreadReply(&GetChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11)})
				return err
			},
			method: http.MethodGet,
			want:   "/api/chat-channels/7/messages/9/thread-replies/11",
		},
		{
			name: "create reply",
			call: func(cl *Client) error {
				_, err := cl.CreateChatChannelThreadReply(&CreateChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReply: reply})
				return err
			},
			method: http.MethodPost,
			want:   "/api/chat-channels/7/messages/9/thread-replies",
			body:   `{"content":"on it","contentType":"TEXT"}`,
		},
		{
			name: "update reply",
			call: func(cl *Client) error {
				_, err := cl.UpdateChatChannelThreadReply(&UpdateChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11), ChatChannelThreadReply: reply})
				return err
			},
			method: http.MethodPut,
			want:   "/api/chat-channels/7/messages/9/thread-replies/11",
			body:   `{"content":"on it","contentType":"TEXT"}`,
		},
		{
			name: "delete reply",
			call: func(cl *Client) error {
				_, err := cl.DeleteChatChannelThreadReply(&DeleteChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11)})
				return err
			},
			method: http.MethodDelete,
			want:   "/api/chat-channels/7/messages/9/thread-replies/11",
		},
		{
			name: "add reaction",
			call: func(cl *Client) error {
				_, err := cl.AddChatChannelThreadReplyReaction(&AddChatChannelThreadReplyReactionInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11), Code: String("eyes")})
				return err
			},
			method: http.MethodPost,
			want:   "/api/chat-channels/7/messages/9/thread-replies/11/reactions",
			body:   `{"code":"eyes"}`,
		},
		{
			name: "remove reaction",
			call: func(cl *Client) error {
				_, err := cl.RemoveChatChannelThreadReplyReaction(&RemoveChatChannelThreadReplyReactionInput{ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11), Code: String("eyes")})
				return err
			},
			method: http.MethodDelete,
			want:   "/api/chat-channels/7/messages/9/thread-replies/11/reactions/eyes",
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
				if r.Method == http.MethodGet && r.URL.Path == "/api/chat-channels/7/messages/9/thread-replies" {
					_, _ = w.Write([]byte(`[]`))
					return
				}
				if r.Method == http.MethodPost && r.URL.Path == "/api/chat-channels/7/messages/9/thread-replies" {
					w.WriteHeader(http.StatusCreated)
				}
				_, _ = w.Write([]byte(`{"id":11,"threadId":9,"content":"on it","contentType":"TEXT"}`))
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
func TestRemoveChatChannelThreadReplyReactionEscapesCode(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":11,"threadId":9,"content":"on it","contentType":"TEXT"}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).RemoveChatChannelThreadReplyReaction(&RemoveChatChannelThreadReplyReactionInput{
		ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11),
		Code: String("thumbs up/down"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/chat-channels/7/messages/9/thread-replies/11/reactions/thumbs%20up%2Fdown" {
		t.Errorf("path = %q, want the code escaped into a single segment", path)
	}
}

// Replies are deleted softly like messages, so the delete answers 200 with the tombstone
// rather than 204.
func TestDeleteChatChannelThreadReplyReturnsTheTombstone(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":11,"threadId":9,"content":"<Deleted on 2026-09-17T11:04:28Z>","contentType":"TEXT","deleted":true}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).DeleteChatChannelThreadReply(&DeleteChatChannelThreadReplyInput{
		ChannelID: Int64(7), ChannelType: String(ChatChannelType.Alert), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.ChatChannelThreadReply.Deleted {
		t.Error("deleted = false, want true")
	}
	if result.ChatChannelThreadReply.ThreadID != 9 {
		t.Errorf("threadId = %d, want 9", result.ChatChannelThreadReply.ThreadID)
	}
}

func TestChatChannelThreadReplyRequiredInputs(t *testing.T) {
	c := NewClient()
	cases := map[string]func() error{
		"nil input": func() error { _, err := c.GetChatChannelThreadReplies(nil); return err },
		"no channel id": func() error {
			_, err := c.GetChatChannelThreadReplies(&GetChatChannelThreadRepliesInput{ChannelType: String("ALERT"), ThreadID: Int64(9)})
			return err
		},
		"no channel type": func() error {
			_, err := c.GetChatChannelThreadReplies(&GetChatChannelThreadRepliesInput{ChannelID: Int64(7), ThreadID: Int64(9)})
			return err
		},
		"no thread id": func() error {
			_, err := c.GetChatChannelThreadReplies(&GetChatChannelThreadRepliesInput{ChannelID: Int64(7), ChannelType: String("ALERT")})
			return err
		},
		"no reply id": func() error {
			_, err := c.GetChatChannelThreadReply(&GetChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ThreadID: Int64(9)})
			return err
		},
		"no reply body": func() error {
			_, err := c.CreateChatChannelThreadReply(&CreateChatChannelThreadReplyInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ThreadID: Int64(9)})
			return err
		},
		"no reaction code": func() error {
			_, err := c.AddChatChannelThreadReplyReaction(&AddChatChannelThreadReplyReactionInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11)})
			return err
		},
		"no code on remove": func() error {
			_, err := c.RemoveChatChannelThreadReplyReaction(&RemoveChatChannelThreadReplyReactionInput{ChannelID: Int64(7), ChannelType: String("ALERT"), ThreadID: Int64(9), ChatChannelThreadReplyID: Int64(11)})
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
