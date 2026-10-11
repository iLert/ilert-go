package ilert

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSupportHourMarshalSupportWindows pins the three things a write has to be able to say about
// the windows: nothing, clear them all, and replace them. A plain slice with omitempty could not
// send the empty list, so the coverage of a support hour could never be cleared through windows.
func TestSupportHourMarshalSupportWindows(t *testing.T) {
	windows := []SupportWindow{{
		From: &TimeOfWeek{DayOfWeek: DayOfWeek.Friday, Time: "17:00"},
		To:   &TimeOfWeek{DayOfWeek: DayOfWeek.Monday, Time: "09:00"},
	}}
	empty := []SupportWindow{}

	cases := []struct {
		name        string
		supportHour *SupportHour
		want        []string
		wantAbsent  []string
	}{
		{
			name:        "nil leaves the windows out",
			supportHour: &SupportHour{Name: "test", SupportDays: &SupportDays{MONDAY: &SupportDay{Start: "09:00", End: "17:00"}}},
			want:        []string{`"supportDays":{"MONDAY":{"start":"09:00","end":"17:00"}`},
			wantAbsent:  []string{`"supportWindows"`},
		},
		{
			name:        "an empty slice clears the coverage",
			supportHour: &SupportHour{Name: "test", SupportWindows: &empty},
			want:        []string{`"supportWindows":[]`, `"supportDays":null`},
		},
		{
			name:        "windows are sent as they are",
			supportHour: &SupportHour{Name: "test", SupportWindows: &windows},
			want:        []string{`"supportWindows":[{"from":{"dayOfWeek":"FRIDAY","time":"17:00"},"to":{"dayOfWeek":"MONDAY","time":"09:00"}}]`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.supportHour)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, want := range tc.want {
				if !strings.Contains(string(payload), want) {
					t.Errorf("payload = %s, want it to contain %s", payload, want)
				}
			}
			for _, absent := range tc.wantAbsent {
				if strings.Contains(string(payload), absent) {
					t.Errorf("payload = %s, want it not to contain %s", payload, absent)
				}
			}
		})
	}
}

// TestSupportHourUnmarshalSupportWindows verifies both response shapes: windows returned next to
// the lossy per-day view, and null windows when the per-day view is the whole coverage.
func TestSupportHourUnmarshalSupportWindows(t *testing.T) {
	t.Run("windows next to the per-day view", func(t *testing.T) {
		payload := []byte(`{"id":1,"name":"test","timezone":"Europe/Berlin",
			"supportDays":{"MONDAY":{"start":"09:00","end":"12:00"},"TUESDAY":null,"WEDNESDAY":null,"THURSDAY":null,"FRIDAY":null,"SATURDAY":null,"SUNDAY":null},
			"supportWindows":[
				{"from":{"dayOfWeek":"MONDAY","time":"09:00"},"to":{"dayOfWeek":"MONDAY","time":"12:00"}},
				{"from":{"dayOfWeek":"MONDAY","time":"13:00"},"to":{"dayOfWeek":"MONDAY","time":"17:00"}}
			]}`)

		var supportHour SupportHour
		if err := json.Unmarshal(payload, &supportHour); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if supportHour.SupportWindows == nil || len(*supportHour.SupportWindows) != 2 {
			t.Fatalf("SupportWindows = %v, want the two windows from the payload", supportHour.SupportWindows)
		}
		second := (*supportHour.SupportWindows)[1]
		if second.From == nil || second.From.DayOfWeek != DayOfWeek.Monday || second.From.Time != "13:00" ||
			second.To == nil || second.To.DayOfWeek != DayOfWeek.Monday || second.To.Time != "17:00" {
			t.Errorf("second window = %+v, want MONDAY 13:00 to MONDAY 17:00", second)
		}
		if supportHour.SupportDays == nil || supportHour.SupportDays.MONDAY == nil || supportHour.SupportDays.MONDAY.End != "12:00" {
			t.Errorf("SupportDays = %+v, want the per-day view of the payload", supportHour.SupportDays)
		}
	})

	t.Run("null windows", func(t *testing.T) {
		payload := []byte(`{"id":1,"name":"test","timezone":"Europe/Berlin",
			"supportDays":{"MONDAY":{"start":"09:00","end":"17:00"}},"supportWindows":null}`)

		var supportHour SupportHour
		if err := json.Unmarshal(payload, &supportHour); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if supportHour.SupportWindows != nil {
			t.Errorf("SupportWindows = %v, want nil", *supportHour.SupportWindows)
		}
	})
}
