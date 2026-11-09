package ilert

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAlertSourceAutoRaiseAlerts(t *testing.T) {
	tests := []struct {
		autoRaiseAlerts *bool
		want            string
	}{
		{nil, ""},
		{Bool(false), `"autoRaiseAlerts":false`},
		{Bool(true), `"autoRaiseAlerts":true`},
	}
	for _, tt := range tests {
		encoded, err := json.Marshal(AlertSource{
			SupportHours:    SupportHoursReference{ID: 7},
			AutoRaiseAlerts: tt.autoRaiseAlerts,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tt.want == "" {
			if strings.Contains(string(encoded), "autoRaiseAlerts") {
				t.Errorf("payload = %s, want no autoRaiseAlerts", encoded)
			}
		} else if !strings.Contains(string(encoded), tt.want) {
			t.Errorf("payload = %s, want it to contain %s", encoded, tt.want)
		}
	}

	alertSource := AlertSource{}
	if err := json.Unmarshal([]byte(`{"supportHours":{"id":7},"autoRaiseAlerts":true}`), &alertSource); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alertSource.AutoRaiseAlerts == nil || !*alertSource.AutoRaiseAlerts {
		t.Errorf("AutoRaiseAlerts = %v, want true", alertSource.AutoRaiseAlerts)
	}
}
