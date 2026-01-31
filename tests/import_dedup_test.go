package tests

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/SebiGabor/personal-finance-cli/internal/cli"
	"github.com/SebiGabor/personal-finance-cli/internal/models"
)

func TestImportDuplicates(t *testing.T) {
	// 1. Setup
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 2. Create a sample CSV file
	csvContent := `Date,Description,Amount,Category
2024-03-01,Test Transaction,50.00,General
`
	tmpfile, err := os.CreateTemp("", "dedup-test-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(csvContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// 3. First Import (Should Succeed)
	out := new(bytes.Buffer)
	cli.RootCmd.SetOut(out)
	cli.RootCmd.SetArgs([]string{"import", tmpfile.Name()})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("First import failed: %v", err)
	}

	// Verify DB has 1 transaction
	txs, _ := models.ListTransactions(db)
	if len(txs) != 1 {
		t.Fatalf("Expected 1 transaction after first import, got %d", len(txs))
	}

	// 4. Second Import (Should Skip)
	out.Reset()
	cli.RootCmd.SetOut(out)
	cli.RootCmd.SetArgs([]string{"import", tmpfile.Name()})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Second import failed: %v", err)
	}

	// 5. Verify Output & DB State
	output := out.String()
	if !strings.Contains(output, "1 duplicates skipped") {
		t.Errorf("Expected output to mention skipped duplicates, got:\n%s", output)
	}

	txs, _ = models.ListTransactions(db)
	if len(txs) != 1 {
		t.Errorf("Expected transaction count to remain 1, but got %d (Duplicates were inserted!)", len(txs))
	}
}
