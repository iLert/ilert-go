package ilert

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// The topology filter is a single condition expression, conditions are combined with "and"
// inside it, so it goes on the query as one value.
func TestGetServiceTopology(t *testing.T) {
	var path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"nodes":[{"id":1,"name":"checkout","status":"OPERATIONAL"},{"id":2,"restricted":true}],"edges":[{"id":9,"sourceServiceId":1,"targetServiceId":2,"type":"HARD"}]}`))
	}))
	defer srv.Close()

	expression := "environment == 'production' and severity != 'low'"
	result, err := newTestClient(t, srv.URL).GetServiceTopology(&GetServiceTopologyInput{Labels: String(expression)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/services/topology" {
		t.Errorf("path = %q, want /api/services/topology", path)
	}
	if got := query.Get("labels"); got != expression {
		t.Errorf("labels = %q, want %q", got, expression)
	}
	if len(query["labels"]) != 1 {
		t.Errorf("labels params = %v, want a single expression", query["labels"])
	}
	if len(result.ServiceGraph.Nodes) != 2 || result.ServiceGraph.Nodes[0].Name != "checkout" {
		t.Errorf("nodes = %+v, want the checkout node first", result.ServiceGraph.Nodes)
	}
	if !result.ServiceGraph.Nodes[1].Restricted {
		t.Errorf("node = %+v, want it marked restricted", result.ServiceGraph.Nodes[1])
	}
	if len(result.ServiceGraph.Edges) != 1 || result.ServiceGraph.Edges[0].TargetServiceID != 2 {
		t.Errorf("edges = %+v, want one edge pointing at service 2", result.ServiceGraph.Edges)
	}
}

func TestGetServiceTopologyWithoutFilter(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"nodes":[],"edges":[]}`))
	}))
	defer srv.Close()

	if _, err := newTestClient(t, srv.URL).GetServiceTopology(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := query["labels"]; ok {
		t.Errorf("labels param is present, want it absent")
	}
}

// The status query param is what keeps a status change that happened in the meantime from
// being published by accident, so it must always reach the API.
func TestPublishServiceStatus(t *testing.T) {
	var method, path string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		// the endpoint answers with the plain service, without the optional publicStatus
		_, _ = w.Write([]byte(`{"id":3,"name":"checkout","status":"MAJOR_OUTAGE"}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).PublishServiceStatus(&PublishServiceStatusInput{
		ServiceID: Int64(3),
		Status:    String(ServiceStatus.MajorOutage),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodPut {
		t.Errorf("method = %s, want PUT", method)
	}
	if path != "/api/services/3/publish-status" {
		t.Errorf("path = %q, want /api/services/3/publish-status", path)
	}
	if got := query.Get("status"); got != "MAJOR_OUTAGE" {
		t.Errorf("status = %q, want MAJOR_OUTAGE", got)
	}
	if result.Service.ID != 3 || result.Service.Status != ServiceStatus.MajorOutage {
		t.Errorf("service = %+v, want service 3 with status MAJOR_OUTAGE", result.Service)
	}
}

func TestPublishServiceStatusRequiredInputs(t *testing.T) {
	c := NewClient()
	if _, err := c.PublishServiceStatus(&PublishServiceStatusInput{Status: String(ServiceStatus.Degraded)}); err == nil {
		t.Error("PublishServiceStatus without a service id = nil error, want an error")
	}
	if _, err := c.PublishServiceStatus(&PublishServiceStatusInput{ServiceID: Int64(3)}); err == nil {
		t.Error("PublishServiceStatus without a status = nil error, want an error")
	}
}
