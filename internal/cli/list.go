package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		transactions, err := models.ListTransactions(database)
		if err != nil {
			return fmt.Errorf("failed to list transactions: %w", err)
		}

		if len(transactions) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No transactions found.")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDATE\tAMOUNT\tCATEGORY\tDESCRIPTION")

		for _, t := range transactions {
			fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\t%s\n",
				t.ID,
				t.Date.Format("2006-01-02"),
				t.Amount,
				t.Category,
				t.Description,
			)
		}
		return w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
