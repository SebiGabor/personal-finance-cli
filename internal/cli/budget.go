package cli

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	budgetCategory string
	budgetAmount   float64
)

// budgetCmd represents the base command for budget management
var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage monthly budgets",
	Long:  "Set, list, and remove monthly spending limits for categories.",
}

var budgetAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Set a budget for a category",
	Example: "finance budget add --category Food --amount 500",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := &models.Budget{
			Category: budgetCategory,
			Amount:   budgetAmount,
			Period:   "monthly",
		}

		if err := models.CreateBudget(database, b); err != nil {
			return fmt.Errorf("failed to create budget: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Budget set: %s -> %.2f/month\n", budgetCategory, budgetAmount)
		return nil
	},
}

var budgetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budgets and current status",
	RunE: func(cmd *cobra.Command, args []string) error {
		budgets, err := models.ListBudgets(database)
		if err != nil {
			return fmt.Errorf("failed to list budgets: %w", err)
		}

		if len(budgets) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No budgets set.")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CATEGORY\tLIMIT\tSPENT\tREMAINING\tSTATUS")

		now := time.Now()

		for _, b := range budgets {
			spent, err := models.GetSpendingTotal(database, b.Category, now.Month(), now.Year())
			if err != nil {
				spent = 0
			}

			remaining := b.Amount - spent
			status := getProgressBar(spent, b.Amount)

			fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%s\n",
				b.Category, b.Amount, spent, remaining, status)
		}
		return w.Flush()
	},
}

var budgetRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a budget by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		if err := models.DeleteBudget(database, id); err != nil {
			return fmt.Errorf("failed to remove budget: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Budget %d removed successfully.\n", id)
		return nil
	},
}

// Helper to create a simple ASCII progress bar
func getProgressBar(spent, limit float64) string {
	if limit == 0 {
		return "[???]"
	}
	percent := spent / limit
	if percent > 1.0 {
		return "[!! OVER BUDGET !!]"
	}

	bars := int(percent * 10) // 0 to 10 bars

	return fmt.Sprintf("[%s%s] %.0f%%",
		strings.Repeat("█", bars),
		strings.Repeat("-", 10-bars),
		percent*100,
	)
}

func init() {
	RootCmd.AddCommand(budgetCmd)
	budgetCmd.AddCommand(budgetAddCmd)
	budgetCmd.AddCommand(budgetListCmd)
	budgetCmd.AddCommand(budgetRemoveCmd)

	budgetAddCmd.Flags().StringVarP(&budgetCategory, "category", "c", "", "Category for the budget")
	budgetAddCmd.Flags().Float64VarP(&budgetAmount, "amount", "a", 0, "Spending limit amount")
	budgetAddCmd.MarkFlagRequired("category")
	budgetAddCmd.MarkFlagRequired("amount")
}
