package cli

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import transactions from a CSV or OFX file",
	Long:  `Imports transactions from a file. Supports .csv and .ofx formats.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		ext := strings.ToLower(filepath.Ext(filePath))

		// Load rules for auto-categorization
		rules, err := models.ListRules(database)
		if err != nil {
			return fmt.Errorf("failed to load categorization rules: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Importing file: %s\n", filePath)

		switch ext {
		case ".csv":
			return importCSV(cmd, filePath, rules)
		case ".ofx":
			return importOFX(cmd, filePath, rules)
		default:
			return fmt.Errorf("unsupported file format '%s'. Please use .csv or .ofx", ext)
		}
	},
}

// --- CSV Logic ---
func importCSV(cmd *cobra.Command, filePath string, rules []models.CategoryRule) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable fields

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %w", err)
	}

	importedCount := 0
	skippedCount := 0

	for i, record := range records {
		if len(record) < 3 {
			continue
		}

		// 1. Date
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

		// 2. Description & Amount
		description := strings.TrimSpace(record[1])
		amountStr := strings.TrimSpace(record[2])
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Row %d skipped: invalid amount\n", i+1)
			skippedCount++
			continue
		}

		// 3. Category
		category := "Uncategorized"
		if len(record) > 3 && strings.TrimSpace(record[3]) != "" {
			category = strings.TrimSpace(record[3])
		} else {
			if match := models.MatchCategory(rules, description); match != "" {
				category = match
			}
		}

		category = models.NormalizeCategory(category)
		tr := &models.Transaction{Date: date, Description: description, Amount: amount, Category: category}
		if err := models.CreateTransaction(database, tr); err != nil {
			skippedCount++
		} else {
			importedCount++
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "CSV Import complete. %d imported, %d skipped.\n", importedCount, skippedCount)
	return nil
}

// --- OFX Logic ---

// XML Structures for parsing OFX
type OFX struct {
	BankMsgs struct {
		StmtTrn struct {
			StmtRs struct {
				BankTranList struct {
					Transactions []OFXTransaction `xml:"STMTTRN"`
				} `xml:"BANKTRANLIST"`
			} `xml:"STMTRS"`
		} `xml:"STMTTRNRS"`
	} `xml:"BANKMSGSRSV1"`
}

type OFXTransaction struct {
	TrnType  string `xml:"TRNTYPE"`
	DtPosted string `xml:"DTPOSTED"`
	TrnAmt   string `xml:"TRNAMT"`
	Name     string `xml:"NAME"`
	Memo     string `xml:"MEMO"`
}

func importOFX(cmd *cobra.Command, filePath string, rules []models.CategoryRule) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse XML
	byteValue, _ := io.ReadAll(file)
	var ofx OFX
	if err := xml.Unmarshal(byteValue, &ofx); err != nil {
		return fmt.Errorf("failed to parse OFX XML: %w", err)
	}

	importedCount := 0
	skippedCount := 0

	for _, t := range ofx.BankMsgs.StmtTrn.StmtRs.BankTranList.Transactions {
		// 1. Parse Date (OFX dates look like "20231025120000" or "20231025")
		date, err := parseOFXDate(t.DtPosted)
		if err != nil {
			skippedCount++
			continue
		}

		// 2. Parse Amount
		amount, err := strconv.ParseFloat(t.TrnAmt, 64)
		if err != nil {
			skippedCount++
			continue
		}

		// 3. Description (Name + Memo)
		description := strings.TrimSpace(t.Name)
		if t.Memo != "" {
			description += " - " + strings.TrimSpace(t.Memo)
		}

		// 4. Auto-Categorize
		category := "Uncategorized"
		if match := models.MatchCategory(rules, description); match != "" {
			category = match
		}

		category = models.NormalizeCategory(category)
		tr := &models.Transaction{Date: date, Description: description, Amount: amount, Category: category}
		if err := models.CreateTransaction(database, tr); err != nil {
			skippedCount++
		} else {
			importedCount++
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "OFX Import complete. %d imported, %d skipped.\n", importedCount, skippedCount)
	return nil
}

func parseOFXDate(dateStr string) (time.Time, error) {
	// Usually YYYYMMDD or YYYYMMDDHHMMSS...
	layout := "20060102"
	if len(dateStr) >= 8 {
		return time.Parse(layout, dateStr[:8])
	}
	return time.Time{}, fmt.Errorf("invalid date length")
}

func init() {
	RootCmd.AddCommand(importCmd)
}
