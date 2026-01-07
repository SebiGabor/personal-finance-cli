package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	reportYear  int
	reportMonth int
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a monthly spending report",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to current date if not provided
		now := time.Now()
		if reportYear == 0 {
			reportYear = now.Year()
		}
		if reportMonth == 0 {
			reportMonth = int(now.Month())
		}

		breakdown, income, expense, err := models.GetMonthlyReport(database, reportYear, reportMonth)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "\n=== Report for %04d-%02d ===\n\n", reportYear, reportMonth)
		fmt.Fprintf(cmd.OutOrStdout(), "Total Income:   %10.2f\n", income)
		fmt.Fprintf(cmd.OutOrStdout(), "Total Expenses: %10.2f\n", expense)
		fmt.Fprintf(cmd.OutOrStdout(), "Net Savings:    %10.2f\n", income+expense)
		fmt.Fprintln(cmd.OutOrStdout(), "\n--- Category Breakdown ---")

		if len(breakdown) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No transactions found for this month.")
			return nil
		}

		// Calculate max absolute amount for scaling the bar chart
		var maxAmount float64
		for _, b := range breakdown {
			if abs(b.Amount) > maxAmount {
				maxAmount = abs(b.Amount)
			}
		}

		// Draw simple ASCII bars
		for _, b := range breakdown {
			barLength := int((abs(b.Amount) / maxAmount) * 20) // Scale to max 20 chars
			bar := strings.Repeat("█", barLength)

			// Adjust spacing for alignment
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s [%-20s] %10.2f\n", b.Category, bar, b.Amount)
		}

		return nil
	},
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func init() {
	RootCmd.AddCommand(reportCmd)
	reportCmd.Flags().IntVarP(&reportYear, "year", "y", 0, "Year of report (default current year)")
	reportCmd.Flags().IntVarP(&reportMonth, "month", "m", 0, "Month of report (default current month)")
}
