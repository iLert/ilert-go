package ilert

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestCallFlowNodeMetadataDefinitions verifies that the metadata of a branch
// node keeps its definitions both ways. Without the field a call flow built from
// CallFlowNodeMetadata could not send a branch node's definitions, and the API
// rejects a DEFINE_BRANCHES node without them. The payload is the shape the API
// returned on 26.10.2026. See iLert/engineering-tasks#2769.
func TestCallFlowNodeMetadataDefinitions(t *testing.T) {
	payload := `{"definitions":[{"branchName":"germany","conditions":"context.fromCountry == 'DE'"}]}`

	var metadata CallFlowNodeMetadata
	if err := json.Unmarshal([]byte(payload), &metadata); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []CallFlowNodeMetadataDefinition{{BranchName: "germany", Conditions: "context.fromCountry == 'DE'"}}
	if !reflect.DeepEqual(metadata.Definitions, want) {
		t.Errorf("Definitions = %v, want %v", metadata.Definitions, want)
	}

	encoded, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != payload {
		t.Errorf("payload = %s, want %s", encoded, payload)
	}
}

// TestEnumsCoverAPIValues verifies that values the API accepts and returns are
// part of the lists the Terraform provider validates against.
func TestEnumsCoverAPIValues(t *testing.T) {
	if !stringSliceContains(CallFlowNodeTypeAll, "DEFINE_BRANCHES") {
		t.Errorf("CallFlowNodeTypeAll = %v, want it to contain DEFINE_BRANCHES", CallFlowNodeTypeAll)
	}
	if !stringSliceContains(CallFlowNodeMetadataAIVoiceModelAll, "ivy") {
		t.Errorf("CallFlowNodeMetadataAIVoiceModelAll = %v, want it to contain ivy", CallFlowNodeMetadataAIVoiceModelAll)
	}
	if !stringSliceContains(StatusPageLayoutAll, "GRID") {
		t.Errorf("StatusPageLayoutAll = %v, want it to contain GRID", StatusPageLayoutAll)
	}
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
