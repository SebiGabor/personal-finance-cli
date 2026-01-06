package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import transactions from a CSV file",
	Long:  `Imports transactions from a CSV file. Expected format: Date,Description,Amount,Category (optional). Date must be YYYY-MM-DD.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read CSV data: %w", err)
		}

		importedCount := 0
		skippedCount := 0

		for i, record := range records {
			// Basic validation: ensure we have at least 3 columns (Date, Desc, Amount)
			if len(record) < 3 {
				continue
			}

			// 1. Parse Date (Assume YYYY-MM-DD for now)
			dateStr := strings.TrimSpace(record[0])
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				// If the first row fails to parse as a date, assume it's a header and skip silently
				if i == 0 {
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Row %d skipped: invalid date format '%s'\n", i+1, dateStr)
				skippedCount++
				continue
			}

			// 2. Parse Description
			description := strings.TrimSpace(record[1])

			// 3. Parse Amount
			amountStr := strings.TrimSpace(record[2])
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Row %d skipped: invalid amount '%s'\n", i+1, amountStr)
				skippedCount++
				continue
			}

			// 4. Parse Category (Optional)
			category := "Uncategorized"
			if len(record) > 3 && strings.TrimSpace(record[3]) != "" {
				category = strings.TrimSpace(record[3])
			}

			// Create and Save Transaction
			tr := &models.Transaction{
				Date:        date,
				Description: description,
				Amount:      amount,
				Category:    category,
			}

			if err := models.CreateTransaction(database, tr); err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Row %d failed to save: %v\n", i+1, err)
				skippedCount++
			} else {
				importedCount++
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Import complete. %d imported, %d skipped.\n", importedCount, skippedCount)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
