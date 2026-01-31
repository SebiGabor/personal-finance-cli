package cli

import (
	"fmt"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new transaction manually",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Get flag values DIRECTLY inside the function
		amount, _ := cmd.Flags().GetFloat64("amount")
		desc, _ := cmd.Flags().GetString("desc")
		catRaw, _ := cmd.Flags().GetString("category")
		dateStr, _ := cmd.Flags().GetString("date")
		category := models.NormalizeCategory(catRaw)

		// 2. Parse Date
		date := time.Now()
		if dateStr != "" {
			var err error
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
			}
		}

		tr := &models.Transaction{
			Date:        date,
			Description: desc,
			Amount:      amount,
			Category:    category,
		}

		// 3. Save Transaction
		if err := models.CreateTransaction(database, tr); err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Successfully added transaction (ID: %d)\n", tr.ID)

		// 4. Budget Alert Logic
		// Only check if it's an expense (negative amount)
		if amount < 0 {
			budgets, _ := models.ListBudgets(database)
			for _, b := range budgets {
				if b.Category == category {
					spent, _ := models.GetSpendingTotal(database, category, date.Month(), date.Year())

					if spent > b.Amount {
						fmt.Fprintf(cmd.OutOrStdout(), "\n⚠️  ALERT: You have exceeded your budget for '%s'!\n", category)
						fmt.Fprintf(cmd.OutOrStdout(), "   Limit: %.2f | Spent: %.2f\n", b.Amount, spent)
					} else if spent > (b.Amount * 0.9) {
						fmt.Fprintf(cmd.OutOrStdout(), "\n⚠️  WARNING: You are close to your budget for '%s' (%.0f%% used).\n", category, (spent/b.Amount)*100)
					}
					break
				}
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// IMPORTANT: Use P-suffix functions (Float64P) and DO NOT pass a pointer (&amount).
	addCmd.Flags().Float64P("amount", "a", 0, "Amount (positive for income, negative for expense)")
	addCmd.Flags().StringP("desc", "d", "", "Transaction description")
	addCmd.Flags().StringP("category", "c", "Uncategorized", "Transaction category")
	addCmd.Flags().StringP("date", "t", "", "Date (YYYY-MM-DD), defaults to today")

	addCmd.MarkFlagRequired("amount")
	addCmd.MarkFlagRequired("desc")
}
