package ilert

import (
	"encoding/json"
	"strings"
	"testing"
)

// A nil Labels map has to disappear from the payload entirely: an absent field leaves
// the labels on the server untouched, while an explicit null would risk clearing them
// on an endpoint whose null handling is not part of its documented contract.
func TestLabelsOmittedWhenNil(t *testing.T) {
	cases := []struct {
		name  string
		value any
	}{
		{"telemetry source", &TelemetrySource{Name: "test", Type: TelemetrySourceType.Otel}},
		{"service", &Service{Name: "test"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Contains(string(payload), `"labels"`) {
				t.Errorf("payload = %s, want no labels field at all", payload)
			}
		})
	}
}

// A pointer to an empty map is the "clear the labels" signal and must reach the API
// as an explicit empty object.
func TestLabelsMarshalEmptyMap(t *testing.T) {
	empty := map[string]string{}
	cases := []struct {
		name  string
		value any
	}{
		{"telemetry source", &TelemetrySource{Name: "test", Type: TelemetrySourceType.Otel, Labels: &empty}},
		{"service", &Service{Name: "test", Labels: &empty}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(string(payload), `"labels":{}`) {
				t.Errorf("payload = %s, want it to contain \"labels\":{}", payload)
			}
		})
	}
}

// Same contract as the other pointer-typed team fields: nil omits, a pointer to an
// empty slice clears.
func TestTelemetrySourceTeamsMarshalling(t *testing.T) {
	payload, err := json.Marshal(&TelemetrySource{Name: "test", Type: TelemetrySourceType.Otel})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(payload), `"teams"`) {
		t.Errorf("payload = %s, want no teams field when Teams is nil", payload)
	}

	empty := []TeamShort{}
	payload, err = json.Marshal(&TelemetrySource{Name: "test", Type: TelemetrySourceType.Otel, Teams: &empty})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(payload), `"teams":[]`) {
		t.Errorf("payload = %s, want it to contain \"teams\":[]", payload)
	}
}

// The type is required on create and every optional field stays out of the payload
// until it is set, so a minimal telemetry source does not overwrite anything the
// account configured elsewhere.
func TestTelemetrySourceMinimalPayload(t *testing.T) {
	payload, err := json.Marshal(&TelemetrySource{Name: "test", Type: TelemetrySourceType.Otel})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["name"] != "test" || decoded["type"] != TelemetrySourceType.Otel {
		t.Errorf("payload = %s, want name and type always sent", payload)
	}
	if len(decoded) != 2 {
		t.Errorf("payload = %s, want only name and type", payload)
	}
}

// The read-only fields the API returns decode into the struct so the provider can
// surface them, the integration key included.
func TestTelemetrySourceUnmarshalReadOnlyFields(t *testing.T) {
	payload := []byte(`{"id":42,"name":"test","type":"OTEL","status":"RECEIVING","integrationKey":"secret","labels":{"env":"prod"},"managedBy":{"type":"TERRAFORM","source":"tf"}}`)

	var source TelemetrySource
	if err := json.Unmarshal(payload, &source); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source.ID != 42 || source.Status != TelemetrySourceStatus.Receiving || source.IntegrationKey != "secret" {
		t.Errorf("source = %+v, want the read-only fields decoded", source)
	}
	if source.Labels == nil || (*source.Labels)["env"] != "prod" {
		t.Errorf("Labels = %v, want the label from the payload", source.Labels)
	}
	if source.ManagedBy == nil || source.ManagedBy.Type != "TERRAFORM" {
		t.Errorf("ManagedBy = %+v, want the payload's managedBy", source.ManagedBy)
	}
}

// Service links follow the same nil-omits / empty-clears contract as labels, and a
// link decodes both of its fields.
func TestServiceLinksMarshalling(t *testing.T) {
	payload, err := json.Marshal(&Service{Name: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(payload), `"links"`) {
		t.Errorf("payload = %s, want no links field when Links is nil", payload)
	}

	empty := []ServiceLink{}
	payload, err = json.Marshal(&Service{Name: "test", Links: &empty})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(payload), `"links":[]`) {
		t.Errorf("payload = %s, want it to contain \"links\":[]", payload)
	}

	var service Service
	if err := json.Unmarshal([]byte(`{"name":"test","links":[{"href":"https://example.com","text":"Runbook"}]}`), &service); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.Links == nil || len(*service.Links) != 1 {
		t.Fatalf("Links = %v, want the single link from the payload", service.Links)
	}
	if (*service.Links)[0].Href != "https://example.com" || (*service.Links)[0].Text != "Runbook" {
		t.Errorf("Links[0] = %+v, want href and text decoded", (*service.Links)[0])
	}
}

// The dependency payload always carries the target service, and leaves the edge id and
// the source service to the endpoint's path when they are unset.
func TestServiceDependencyMarshalling(t *testing.T) {
	payload, err := json.Marshal(&ServiceDependency{TargetServiceID: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["targetServiceId"] != float64(7) {
		t.Errorf("targetServiceId = %v, want 7", decoded["targetServiceId"])
	}
	for _, field := range []string{"id", "sourceServiceId", "notes", "invalidAfter", "createdAt", "updatedAt"} {
		if _, ok := decoded[field]; ok {
			t.Errorf("payload = %s, want %s omitted when unset", payload, field)
		}
	}
}
