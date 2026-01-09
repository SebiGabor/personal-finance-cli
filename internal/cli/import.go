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
	Long:  `Imports transactions from a CSV file. Expected format: Date,Description,Amount,Category (optional).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Load rules first to use for auto-categorization
		rules, err := models.ListRules(database)
		if err != nil {
			return fmt.Errorf("failed to load categorization rules: %w", err)
		}

		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		// ALLOW VARIABLE NUMBER OF FIELDS
		// This prevents the "wrong number of fields" error if the Category column is missing
		reader.FieldsPerRecord = -1

		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read CSV data: %w", err)
		}

		importedCount := 0
		skippedCount := 0

		for i, record := range records {
			// Ensure we have at least 3 columns (Date, Desc, Amount)
			if len(record) < 3 {
				continue
			}

			// Parse Date
			dateStr := strings.TrimSpace(record[0])
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				if i == 0 {
					continue
				} // Skip header
				fmt.Fprintf(cmd.OutOrStdout(), "Row %d skipped: invalid date\n", i+1)
				skippedCount++
				continue
			}

			description := strings.TrimSpace(record[1])

			// Parse Amount
			amountStr := strings.TrimSpace(record[2])
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Row %d skipped: invalid amount\n", i+1)
				skippedCount++
				continue
			}

			// Determine Category
			category := "Uncategorized"
			// If CSV has a category column, use it
			if len(record) > 3 && strings.TrimSpace(record[3]) != "" {
				category = strings.TrimSpace(record[3])
			} else {
				// Otherwise, try to auto-categorize
				if match := models.MatchCategory(rules, description); match != "" {
					category = match
				}
			}

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
