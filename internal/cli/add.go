package cli

import (
	"fmt"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	amount   float64
	desc     string
	category string
	dateStr  string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new transaction manually",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		err := models.CreateTransaction(database, tr)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Successfully added transaction (ID: %d)\n", tr.ID)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Amount (positive for income, negative for expense)")
	addCmd.Flags().StringVarP(&desc, "desc", "d", "", "Transaction description")
	addCmd.Flags().StringVarP(&category, "category", "c", "Uncategorized", "Transaction category")
	addCmd.Flags().StringVarP(&dateStr, "date", "t", "", "Date (YYYY-MM-DD), defaults to today")

	addCmd.MarkFlagRequired("amount")
	addCmd.MarkFlagRequired("desc")
}
