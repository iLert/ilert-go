package ilert

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestCreateUserQueryParams verifies which query parameters CreateUser puts on the
// request. purchase-seat buys a license and charges the account for it, so it must
// never reach the API unless the caller set it explicitly.
// See iLert/engineering-tasks#2391.
func TestCreateUserQueryParams(t *testing.T) {
	cases := []struct {
		name   string
		input  *CreateUserInput
		want   map[string]string
		absent []string
	}{
		{
			name:   "no options",
			input:  &CreateUserInput{User: &User{Email: "test@example.com"}},
			want:   map[string]string{},
			absent: []string{"send-no-invitation", "purchase-seat"},
		},
		{
			name:   "purchase seat",
			input:  &CreateUserInput{User: &User{Email: "test@example.com"}, PurchaseSeat: Bool(true)},
			want:   map[string]string{"purchase-seat": "true"},
			absent: []string{"send-no-invitation"},
		},
		{
			name:   "purchase seat disabled",
			input:  &CreateUserInput{User: &User{Email: "test@example.com"}, PurchaseSeat: Bool(false)},
			want:   map[string]string{"purchase-seat": "false"},
			absent: []string{"send-no-invitation"},
		},
		{
			name: "both options",
			input: &CreateUserInput{
				User:             &User{Email: "test@example.com"},
				SendNoInvitation: Bool(true),
				PurchaseSeat:     Bool(true),
			},
			want: map[string]string{"send-no-invitation": "true", "purchase-seat": "true"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var query url.Values
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query = r.URL.Query()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":1,"email":"test@example.com"}`))
			}))
			defer srv.Close()

			if _, err := newTestClient(t, srv.URL).CreateUser(tc.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for key, want := range tc.want {
				if got := query.Get(key); got != want {
					t.Errorf("query %s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.absent {
				if _, ok := query[key]; ok {
					t.Errorf("query contains %s = %q, want it to be absent", key, query.Get(key))
				}
			}
		})
	}
}
