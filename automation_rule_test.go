package ilert

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAutomationRuleWithLabelledService(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"0a2968c7","alertType":"CREATED","resolveService":true,"serviceStatus":"DEGRADED","service":{"id":7,"name":"checkout","labels":{"entry":[{"key":"team","value":"payments"}]}},"alertSource":{"id":3}}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetAutomationRule(&GetAutomationRuleInput{AutomationRuleID: String("0a2968c7")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/automation-rules/0a2968c7" {
		t.Errorf("path = %s, want /api/automation-rules/0a2968c7", path)
	}
	service := result.AutomationRule.Service
	if service == nil || service.ID != 7 || service.Labels == nil {
		t.Fatalf("service = %+v, want service 7 with its labels", service)
	}
	if labels := *service.Labels; len(labels) != 1 || labels["team"] != "payments" {
		t.Errorf("labels = %v, want team=payments", labels)
	}
}
