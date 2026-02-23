package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicolasacchi/ga-cli/internal/auth"
	"github.com/nicolasacchi/ga-cli/internal/output"
	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var (
	realtimeDimensions string
	realtimeMetrics    string
	realtimeLimit      int64
)

var realtimeCmd = &cobra.Command{
	Use:   "realtime <property-id>",
	Short: "Run a realtime report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := auth.NewAnalyticsDataService(ctx)
		if err != nil {
			return err
		}

		propertyID := args[0]

		req := &analyticsdata.RunRealtimeReportRequest{
			Limit: realtimeLimit,
		}

		if realtimeMetrics != "" {
			for _, m := range strings.Split(realtimeMetrics, ",") {
				req.Metrics = append(req.Metrics, &analyticsdata.Metric{Name: strings.TrimSpace(m)})
			}
		} else {
			req.Metrics = []*analyticsdata.Metric{
				{Name: "activeUsers"},
			}
		}

		if realtimeDimensions != "" {
			for _, d := range strings.Split(realtimeDimensions, ",") {
				req.Dimensions = append(req.Dimensions, &analyticsdata.Dimension{Name: strings.TrimSpace(d)})
			}
		}

		name := fmt.Sprintf("properties/%s", propertyID)
		resp, err := svc.Properties.RunRealtimeReport(name, req).Do()
		if err != nil {
			return err
		}

		return output.PrintJSON(flattenRealtimeReport(resp))
	},
}

func flattenRealtimeReport(resp *analyticsdata.RunRealtimeReportResponse) map[string]any {
	var dimHeaders []string
	for _, h := range resp.DimensionHeaders {
		dimHeaders = append(dimHeaders, h.Name)
	}
	var metricHeaders []string
	for _, h := range resp.MetricHeaders {
		metricHeaders = append(metricHeaders, h.Name)
	}

	var rows []map[string]string
	for _, row := range resp.Rows {
		r := make(map[string]string)
		for i, dv := range row.DimensionValues {
			if i < len(dimHeaders) {
				r[dimHeaders[i]] = dv.Value
			}
		}
		for i, mv := range row.MetricValues {
			if i < len(metricHeaders) {
				r[metricHeaders[i]] = mv.Value
			}
		}
		rows = append(rows, r)
	}

	return map[string]any{
		"row_count": resp.RowCount,
		"rows":      rows,
	}
}

func init() {
	realtimeCmd.Flags().StringVar(&realtimeDimensions, "dimensions", "", "Comma-separated dimensions")
	realtimeCmd.Flags().StringVar(&realtimeMetrics, "metrics", "", "Comma-separated metrics (default: activeUsers)")
	realtimeCmd.Flags().Int64Var(&realtimeLimit, "limit", 100, "Max rows to return")
}
