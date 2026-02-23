package commands

import (
	"context"
	"fmt"

	"github.com/nicolasacchi/ga-cli/internal/auth"
	"github.com/nicolasacchi/ga-cli/internal/output"
	"github.com/spf13/cobra"
)

var propertiesCmd = &cobra.Command{
	Use:   "properties",
	Short: "Manage GA4 properties",
}

var propertiesGetCmd = &cobra.Command{
	Use:   "get <property-id>",
	Short: "Get property details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := auth.NewAnalyticsAdminService(ctx)
		if err != nil {
			return err
		}

		name := fmt.Sprintf("properties/%s", args[0])
		prop, err := svc.Properties.Get(name).Do()
		if err != nil {
			return err
		}

		result := map[string]any{
			"name":         prop.Name,
			"displayName":  prop.DisplayName,
			"propertyType": prop.PropertyType,
			"timeZone":     prop.TimeZone,
			"currencyCode": prop.CurrencyCode,
			"industryCategory": prop.IndustryCategory,
			"createTime":   prop.CreateTime,
			"updateTime":   prop.UpdateTime,
		}
		if prop.Parent != "" {
			result["parent"] = prop.Parent
		}

		return output.PrintJSON(result)
	},
}

func init() {
	propertiesCmd.AddCommand(propertiesGetCmd)
}
