package ilert

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetCallFlowNumbersQueryParams(t *testing.T) {
	cases := []struct {
		name   string
		input  *GetCallFlowNumbersInput
		want   map[string]string
		absent []string
	}{
		{
			name:   "no options",
			input:  nil,
			want:   map[string]string{},
			absent: []string{"state", "query", "include", "start-index", "max-results"},
		},
		{
			name: "available numbers with assignment",
			input: &GetCallFlowNumbersInput{
				State:      String(CallFlowNumberStateFilter.Available),
				Include:    []*string{String(CallFlowNumberInclude.AssignedTo)},
				Query:      String("support"),
				StartIndex: Int(0),
				MaxResults: Int(100),
			},
			want: map[string]string{
				"state":       "AVAILABLE",
				"include":     "assignedTo",
				"query":       "support",
				"start-index": "0",
				"max-results": "100",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var query url.Values
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query = r.URL.Query()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"id":1,"name":"support hotline","state":"USED","phoneNumber":{"regionCode":"DE","number":"+4930123456"},"assignedTo":{"id":7,"name":"support call flow"}}]`))
			}))
			defer srv.Close()

			result, err := newTestClient(t, srv.URL).GetCallFlowNumbers(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for key, want := range tc.want {
				if got := query.Get(key); got != want {
					t.Errorf("query %s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.absent {
				if _, ok := query[key]; ok {
					t.Errorf("query %s is present, want it absent", key)
				}
			}
			if len(result.CallFlowNumbers) != 1 {
				t.Fatalf("call flow numbers = %+v, want one number", result.CallFlowNumbers)
			}
			number := result.CallFlowNumbers[0]
			if number.State != CallFlowNumberState.Used {
				t.Errorf("state = %q, want USED", number.State)
			}
			if number.AssignedTo == nil || number.AssignedTo.ID != 7 {
				t.Errorf("assignedTo = %+v, want the call flow with id 7", number.AssignedTo)
			}
			if number.PhoneNumber == nil || number.PhoneNumber.Number != "+4930123456" {
				t.Errorf("phoneNumber = %+v, want +4930123456", number.PhoneNumber)
			}
		})
	}
}

// A call flow number name is user input and sits in the path of the search endpoint.
func TestSearchCallFlowNumberEscapesName(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"name":"support hotline"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).SearchCallFlowNumber(&SearchCallFlowNumberInput{CallFlowNumberName: String("support hotline")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/call-flow-numbers/name/support%20hotline" {
		t.Errorf("path = %q, want the name escaped into a single segment", path)
	}
	if result.CallFlowNumber.Name != "support hotline" {
		t.Errorf("name = %q, want \"support hotline\"", result.CallFlowNumber.Name)
	}
}

func TestCallFlowNumberRequiredInputs(t *testing.T) {
	c := NewClient()
	if _, err := c.GetCallFlowNumber(&GetCallFlowNumberInput{}); err == nil {
		t.Error("GetCallFlowNumber without an id = nil error, want an error")
	}
	if _, err := c.SearchCallFlowNumber(&SearchCallFlowNumberInput{}); err == nil {
		t.Error("SearchCallFlowNumber without a name = nil error, want an error")
	}
}
