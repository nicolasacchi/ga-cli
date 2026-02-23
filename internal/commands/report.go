package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nicolasacchi/ga-cli/internal/auth"
	"github.com/nicolasacchi/ga-cli/internal/output"
	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var (
	reportStartDate  string
	reportEndDate    string
	reportDimensions string
	reportMetrics    string
	reportLimit      int64
	reportOffset     int64
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Run GA4 reports",
}

var reportRunCmd = &cobra.Command{
	Use:   "run <property-id>",
	Short: "Run a standard GA4 report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := auth.NewAnalyticsDataService(ctx)
		if err != nil {
			return err
		}

		propertyID := args[0]

		startDate := reportStartDate
		if startDate == "" {
			startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		}
		endDate := reportEndDate
		if endDate == "" {
			endDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		}

		req := &analyticsdata.RunReportRequest{
			DateRanges: []*analyticsdata.DateRange{
				{StartDate: startDate, EndDate: endDate},
			},
			Limit:  reportLimit,
			Offset: reportOffset,
		}

		if reportMetrics != "" {
			for _, m := range strings.Split(reportMetrics, ",") {
				req.Metrics = append(req.Metrics, &analyticsdata.Metric{Name: strings.TrimSpace(m)})
			}
		} else {
			req.Metrics = []*analyticsdata.Metric{
				{Name: "sessions"},
				{Name: "activeUsers"},
			}
		}

		if reportDimensions != "" {
			for _, d := range strings.Split(reportDimensions, ",") {
				req.Dimensions = append(req.Dimensions, &analyticsdata.Dimension{Name: strings.TrimSpace(d)})
			}
		}

		name := fmt.Sprintf("properties/%s", propertyID)
		resp, err := svc.Properties.RunReport(name, req).Do()
		if err != nil {
			return err
		}

		return output.PrintJSON(flattenReport(resp))
	},
}

func flattenReport(resp *analyticsdata.RunReportResponse) map[string]any {
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
	reportRunCmd.Flags().StringVar(&reportStartDate, "start-date", "", "Start date (YYYY-MM-DD, default: 7 days ago)")
	reportRunCmd.Flags().StringVar(&reportEndDate, "end-date", "", "End date (YYYY-MM-DD, default: yesterday)")
	reportRunCmd.Flags().StringVar(&reportDimensions, "dimensions", "", "Comma-separated dimensions (e.g. date,country)")
	reportRunCmd.Flags().StringVar(&reportMetrics, "metrics", "", "Comma-separated metrics (e.g. sessions,activeUsers)")
	reportRunCmd.Flags().Int64Var(&reportLimit, "limit", 100, "Max rows to return")
	reportRunCmd.Flags().Int64Var(&reportOffset, "offset", 0, "Row offset for pagination")
	reportCmd.AddCommand(reportRunCmd)
}
