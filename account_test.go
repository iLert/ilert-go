package ilert

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentAccount(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"acc-1","organizationName":"ilert","timezone":"Europe/Berlin","language":"en","region":"DE","enforceMobileProtection":true,"allowAdminSeatPurchase":true,"aiMode":"EU","applicationFeatures":["CALL_ROUTING"],"subscription":{"name":"Premium","status":"ACTIVE"}}`))
	}))
	defer srv.Close()

	result, err := newTestClient(t, srv.URL).GetCurrentAccount(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/api/tenants/current" {
		t.Errorf("path = %q, want /api/tenants/current", path)
	}
	account := result.Account
	if account.ID != "acc-1" || account.OrganizationName != "ilert" {
		t.Errorf("account = %+v, want the ilert account acc-1", account)
	}
	if account.AiMode != AiMode.EU {
		t.Errorf("aiMode = %q, want EU", account.AiMode)
	}
	if !account.AllowAdminSeatPurchase {
		t.Error("allowAdminSeatPurchase = false, want true")
	}
	if account.Subscription == nil || account.Subscription.Status != SubscriptionStatus.Active {
		t.Errorf("subscription = %+v, want an active plan", account.Subscription)
	}
	if len(account.ApplicationFeatures) != 1 || account.ApplicationFeatures[0] != "CALL_ROUTING" {
		t.Errorf("applicationFeatures = %v, want [CALL_ROUTING]", account.ApplicationFeatures)
	}
}
