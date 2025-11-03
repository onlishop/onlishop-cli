package account

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	account_api "github.com/onlishop/onlishop-cli/internal/account-api"
	"github.com/onlishop/onlishop-cli/internal/table"
)

var accountCompanyProducerExtensionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all your extensions",
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := services.AccountClient.Producer()
		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		criteria := account_api.ListExtensionCriteria{
			Limit: 100,
		}

		if len(listExtensionSearch) > 0 {
			criteria.Query = &account_api.Query{
				Type:  "equals",
				Field: "productNumber",
				Value: listExtensionSearch,
			}
			criteria.OrderBy = "created_at"
			criteria.OrderSequence = "desc"
		}

		extensions, err := p.Extensions(cmd.Context(), &criteria)
		if err != nil {
			return err
		}

		table := table.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Name", "Type", "Compatible with latest version", "Status"})

		for _, extension := range extensions {
			if extension.Status.Name == "deleted" {
				continue
			}

			compatible := "No"

			if extension.IsCompatibleWithLatestOnlishopVersion {
				compatible = "Yes"
			}

			_ = table.Append([]string{
				extension.Id,
				extension.Name,
				extension.Generation.Description,
				compatible,
				extension.Status.Name,
			})
		}

		_ = table.Render()

		return nil
	},
}

var listExtensionSearch string

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionListCmd)
	accountCompanyProducerExtensionListCmd.Flags().StringVar(&listExtensionSearch, "search", "", "Filter for name")
}
