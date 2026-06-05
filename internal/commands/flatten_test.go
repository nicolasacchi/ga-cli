package commands

import (
	"testing"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

func TestFlattenReport(t *testing.T) {
	resp := &analyticsdata.RunReportResponse{
		DimensionHeaders: []*analyticsdata.DimensionHeader{{Name: "date"}, {Name: "country"}},
		MetricHeaders:    []*analyticsdata.MetricHeader{{Name: "sessions"}, {Name: "activeUsers"}},
		Rows: []*analyticsdata.Row{
			{
				DimensionValues: []*analyticsdata.DimensionValue{{Value: "20260222"}, {Value: "IT"}},
				MetricValues:    []*analyticsdata.MetricValue{{Value: "1234"}, {Value: "987"}},
			},
		},
		RowCount: 1,
	}

	out := flattenReport(resp)
	if out["row_count"] != int64(1) {
		t.Errorf("row_count = %v, want 1", out["row_count"])
	}
	rows, ok := out["rows"].([]map[string]string)
	if !ok || len(rows) != 1 {
		t.Fatalf("rows = %#v", out["rows"])
	}
	r := rows[0]
	for k, want := range map[string]string{"date": "20260222", "country": "IT", "sessions": "1234", "activeUsers": "987"} {
		if r[k] != want {
			t.Errorf("row[%q] = %q, want %q", k, r[k], want)
		}
	}
}

// TestFlattenReport_MismatchedValuesDoNotPanic: more values than headers must be
// ignored by the i<len(headers) guard rather than panicking.
func TestFlattenReport_MismatchedValuesDoNotPanic(t *testing.T) {
	resp := &analyticsdata.RunReportResponse{
		DimensionHeaders: []*analyticsdata.DimensionHeader{{Name: "date"}},
		MetricHeaders:    []*analyticsdata.MetricHeader{{Name: "sessions"}},
		Rows: []*analyticsdata.Row{
			{
				DimensionValues: []*analyticsdata.DimensionValue{{Value: "20260222"}, {Value: "extra"}},
				MetricValues:    []*analyticsdata.MetricValue{{Value: "1234"}, {Value: "extra"}},
			},
		},
		RowCount: 1,
	}
	out := flattenReport(resp)
	rows := out["rows"].([]map[string]string)
	if len(rows[0]) != 2 {
		t.Errorf("only headered values should map, got %#v", rows[0])
	}
}

func TestFlattenReport_Empty(t *testing.T) {
	out := flattenReport(&analyticsdata.RunReportResponse{})
	if out["row_count"] != int64(0) {
		t.Errorf("empty row_count = %v, want 0", out["row_count"])
	}
	if rows := out["rows"].([]map[string]string); len(rows) != 0 {
		t.Errorf("empty rows = %#v", rows)
	}
}

func TestFlattenRealtimeReport(t *testing.T) {
	resp := &analyticsdata.RunRealtimeReportResponse{
		DimensionHeaders: []*analyticsdata.DimensionHeader{{Name: "country"}},
		MetricHeaders:    []*analyticsdata.MetricHeader{{Name: "activeUsers"}},
		Rows: []*analyticsdata.Row{
			{
				DimensionValues: []*analyticsdata.DimensionValue{{Value: "IT"}},
				MetricValues:    []*analyticsdata.MetricValue{{Value: "42"}},
			},
		},
		RowCount: 1,
	}
	out := flattenRealtimeReport(resp)
	rows := out["rows"].([]map[string]string)
	if len(rows) != 1 || rows[0]["country"] != "IT" || rows[0]["activeUsers"] != "42" {
		t.Errorf("realtime flatten wrong: %#v", rows)
	}
}
