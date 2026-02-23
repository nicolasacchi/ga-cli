package commands

import (
	"context"

	"github.com/nicolasacchi/ga-cli/internal/auth"
	"github.com/nicolasacchi/ga-cli/internal/output"
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage GA4 accounts",
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all GA4 account and property summaries",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := auth.NewAnalyticsAdminService(ctx)
		if err != nil {
			return err
		}

		resp, err := svc.AccountSummaries.List().Do()
		if err != nil {
			return err
		}

		type propertySummary struct {
			Property    string `json:"property"`
			DisplayName string `json:"displayName"`
		}
		type accountSummary struct {
			Account     string            `json:"account"`
			DisplayName string            `json:"displayName"`
			Properties  []propertySummary `json:"properties"`
		}

		var result []accountSummary
		for _, a := range resp.AccountSummaries {
			as := accountSummary{
				Account:     a.Account,
				DisplayName: a.DisplayName,
			}
			for _, p := range a.PropertySummaries {
				as.Properties = append(as.Properties, propertySummary{
					Property:    p.Property,
					DisplayName: p.DisplayName,
				})
			}
			result = append(result, as)
		}

		return output.PrintJSON(result)
	},
}

func init() {
	accountsCmd.AddCommand(accountsListCmd)
}
