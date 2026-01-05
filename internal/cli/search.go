package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search transactions by description or category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		transactions, err := models.SearchTransactions(database, query)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(transactions) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "No transactions found matching '%s'.\n", query)
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDATE\tAMOUNT\tCATEGORY\tDESCRIPTION")
		for _, t := range transactions {
			fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\t%s\n",
				t.ID, t.Date.Format("2006-01-02"), t.Amount, t.Category, t.Description)
		}
		return w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
