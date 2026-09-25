package ilert

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Nothing the caller did not set may end up on the query: the alert label endpoints reject
// any page size other than their fixed one, so an invented default would 400 every plain
// call. The service and telemetry source endpoints page freely.
func TestLabelKeysQueryParams(t *testing.T) {
	cases := []struct {
		name   string
		input  *GetServiceLabelKeysInput
		want   map[string]string
		absent []string
	}{
		{
			name:   "no options",
			input:  nil,
			want:   map[string]string{},
			absent: []string{"start-index", "max-results", "query"},
		},
		{
			name:  "second page with prefix",
			input: &GetServiceLabelKeysInput{StartIndex: Int(100), MaxResults: Int(100), Query: String("env")},
			want:  map[string]string{"max-results": "100", "start-index": "100", "query": "env"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var path string
			var query url.Values
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				query = r.URL.Query()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"key":"environment"}]`))
			}))
			defer srv.Close()

			result, err := newTestClient(t, srv.URL).GetServiceLabelKeys(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if path != "/api/services/labels" {
				t.Errorf("path = %q, want /api/services/labels", path)
			}
			for key, want := range tc.want {
				if _, ok := query[key]; !ok {
					t.Errorf("query %s is absent, want %q", key, want)
					continue
				}
				if got := query.Get(key); got != want {
					t.Errorf("query %s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.absent {
				if _, ok := query[key]; ok {
					t.Errorf("query %s is present, want it absent", key)
				}
			}
			if len(result.LabelKeys) != 1 || result.LabelKeys[0].Key != "environment" {
				t.Errorf("label keys = %+v, want one key \"environment\"", result.LabelKeys)
			}
		})
	}
}

// Each entity has its own label namespace, so the three pairs of methods must hit the
// route of the entity they are named after.
func TestLabelRoutesPerEntity(t *testing.T) {
	cases := []struct {
		name string
		call func(*Client) error
		want string
	}{
		{
			name: "alert keys",
			call: func(c *Client) error { _, err := c.GetAlertLabelKeys(nil); return err },
			want: "/api/alerts/labels",
		},
		{
			name: "service keys",
			call: func(c *Client) error { _, err := c.GetServiceLabelKeys(nil); return err },
			want: "/api/services/labels",
		},
		{
			name: "telemetry source keys",
			call: func(c *Client) error { _, err := c.GetTelemetrySourceLabelKeys(nil); return err },
			want: "/api/telemetry-sources/labels",
		},
		{
			name: "alert values",
			call: func(c *Client) error {
				_, err := c.GetAlertLabelValues(&GetAlertLabelValuesInput{LabelKey: String("environment")})
				return err
			},
			want: "/api/alerts/labels/environment/values",
		},
		{
			name: "service values",
			call: func(c *Client) error {
				_, err := c.GetServiceLabelValues(&GetServiceLabelValuesInput{LabelKey: String("environment")})
				return err
			},
			want: "/api/services/labels/environment/values",
		},
		{
			name: "telemetry source values",
			call: func(c *Client) error {
				_, err := c.GetTelemetrySourceLabelValues(&GetTelemetrySourceLabelValuesInput{LabelKey: String("environment")})
				return err
			},
			want: "/api/telemetry-sources/labels/environment/values",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var path string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[]`))
			}))
			defer srv.Close()

			if err := tc.call(newTestClient(t, srv.URL)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if path != tc.want {
				t.Errorf("path = %q, want %q", path, tc.want)
			}
		})
	}
}

// A label key is free-form user input and sits in the path, so it has to be escaped
// instead of splitting the route into extra segments.
func TestLabelValuesEscapesKey(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"value":"production"}]`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetServiceLabelValues(&GetServiceLabelValuesInput{LabelKey: String("ilert.com/discovered-by")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/services/labels/ilert.com%2Fdiscovered-by/values" {
		t.Errorf("path = %q, want the key escaped into a single segment", path)
	}
	if len(result.LabelValues) != 1 || result.LabelValues[0].Value != "production" {
		t.Errorf("label values = %+v, want one value \"production\"", result.LabelValues)
	}
}

// The alert label endpoints reject a page size that is not AlertLabelPageSize, so the
// values path must not invent one either.
func TestLabelValuesSendsNoUnsetParams(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	if _, err := newTestClient(t, srv.URL).GetAlertLabelValues(&GetAlertLabelValuesInput{LabelKey: String("environment")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"query", "start-index", "max-results"} {
		if _, ok := query[key]; ok {
			t.Errorf("query %s is present, want it absent", key)
		}
	}
}

// Paging the alert label keys means stepping in AlertLabelPageSize, which is the only
// page size that endpoint accepts.
func TestAlertLabelKeysPaging(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv.URL).GetAlertLabelKeys(&GetAlertLabelKeysInput{
		StartIndex: Int(AlertLabelPageSize),
		MaxResults: Int(AlertLabelPageSize),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := query.Get("start-index"); got != "100" {
		t.Errorf("start-index = %q, want 100", got)
	}
	if got := query.Get("max-results"); got != "100" {
		t.Errorf("max-results = %q, want 100", got)
	}
}

// The label key sits in the path, so a missing one has to fail before the request is
// built rather than producing a request against "/labels//values". Each entity has its
// own input type, so each one needs the guard.
func TestLabelValuesRequiresKey(t *testing.T) {
	c := NewClient()
	cases := []struct {
		name     string
		nilInput func() error
		noKey    func() error
	}{
		{
			name:     "alerts",
			nilInput: func() error { _, err := c.GetAlertLabelValues(nil); return err },
			noKey:    func() error { _, err := c.GetAlertLabelValues(&GetAlertLabelValuesInput{}); return err },
		},
		{
			name:     "services",
			nilInput: func() error { _, err := c.GetServiceLabelValues(nil); return err },
			noKey:    func() error { _, err := c.GetServiceLabelValues(&GetServiceLabelValuesInput{}); return err },
		},
		{
			name:     "telemetry sources",
			nilInput: func() error { _, err := c.GetTelemetrySourceLabelValues(nil); return err },
			noKey: func() error {
				_, err := c.GetTelemetrySourceLabelValues(&GetTelemetrySourceLabelValuesInput{})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.nilInput(); err == nil {
				t.Error("nil input = nil error, want an error")
			}
			if err := tc.noKey(); err == nil {
				t.Error("input without a key = nil error, want an error")
			}
		})
	}
}

// A nil input is the "no options" case on the key endpoints, not an error, matching
// GetNumbers and the other listings that take no required field.
func TestLabelKeysAcceptNilInput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"key":"environment"}]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	cases := map[string]func() error{
		"alerts":            func() error { _, err := c.GetAlertLabelKeys(nil); return err },
		"services":          func() error { _, err := c.GetServiceLabelKeys(nil); return err },
		"telemetry sources": func() error { _, err := c.GetTelemetrySourceLabelKeys(nil); return err },
	}

	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			if err := call(); err != nil {
				t.Errorf("nil input = %v, want no error", err)
			}
		})
	}
}

// The alert label endpoints answer guests and stakeholders with an empty object rather than
// an empty list. That has to read as no labels, while any other object still fails to decode.
func TestAlertLabelsReadEmptyObjectAsNoLabels(t *testing.T) {
	emptyObject := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer emptyObject.Close()

	c := newTestClient(t, emptyObject.URL)
	keys, err := c.GetAlertLabelKeys(nil)
	if err != nil {
		t.Fatalf("GetAlertLabelKeys() error = %v, want none", err)
	}
	if keys.LabelKeys == nil || len(keys.LabelKeys) != 0 {
		t.Errorf("label keys = %v, want an empty list", keys.LabelKeys)
	}
	values, err := c.GetAlertLabelValues(&GetAlertLabelValuesInput{LabelKey: String("environment")})
	if err != nil {
		t.Fatalf("GetAlertLabelValues() error = %v, want none", err)
	}
	if values.LabelValues == nil || len(values.LabelValues) != 0 {
		t.Errorf("label values = %v, want an empty list", values.LabelValues)
	}

	otherObject := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"key":"environment"}`))
	}))
	defer otherObject.Close()

	if _, err := newTestClient(t, otherObject.URL).GetAlertLabelKeys(nil); err == nil {
		t.Error("GetAlertLabelKeys() on an object that is not empty = nil error, want a decoding error")
	}
}
