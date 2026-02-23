package commands

import (
	"context"
	"fmt"

	"github.com/nicolasacchi/ga-cli/internal/auth"
	"github.com/nicolasacchi/ga-cli/internal/output"
	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata <property-id>",
	Short: "List available dimensions and metrics for a property",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := auth.NewAnalyticsDataService(ctx)
		if err != nil {
			return err
		}

		name := fmt.Sprintf("properties/%s/metadata", args[0])
		resp, err := svc.Properties.GetMetadata(name).Do()
		if err != nil {
			return err
		}

		type dimInfo struct {
			APIName     string `json:"apiName"`
			UIName      string `json:"uiName"`
			Description string `json:"description"`
			Category    string `json:"category"`
		}
		type metricInfo struct {
			APIName     string `json:"apiName"`
			UIName      string `json:"uiName"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Type        string `json:"type"`
		}

		var dims []dimInfo
		for _, d := range resp.Dimensions {
			dims = append(dims, dimInfo{
				APIName:     d.ApiName,
				UIName:      d.UiName,
				Description: d.Description,
				Category:    d.Category,
			})
		}

		var metrics []metricInfo
		for _, m := range resp.Metrics {
			metrics = append(metrics, metricInfo{
				APIName:     m.ApiName,
				UIName:      m.UiName,
				Description: m.Description,
				Category:    m.Category,
				Type:        m.Type,
			})
		}

		return output.PrintJSON(map[string]any{
			"dimensions": dims,
			"metrics":    metrics,
		})
	},
}
