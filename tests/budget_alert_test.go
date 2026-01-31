package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/SebiGabor/personal-finance-cli/internal/cli"
	"github.com/spf13/pflag"
)

func TestBudgetAlerts(t *testing.T) {
	// 1. Initialize DB
	db := NewTestDB(t)
	cli.SetDatabase(db)

	for _, cmd := range cli.RootCmd.Commands() {
		if cmd.Name() == "add" {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if err := f.Value.Set(f.DefValue); err != nil {
					t.Logf("Failed to reset flag %s: %v", f.Name, err)
				}
				f.Changed = false
			})
		}
	}

	// 2. Set a Budget of 100 for "Entertainment"
	cli.RootCmd.SetArgs([]string{"budget", "add", "--category", "Entertainment", "--amount", "100"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Budget add failed: %v", err)
	}

	// 3. Spend 50 (50% used)
	// We use --amount=-50 to ensure it's parsed as a negative value
	cli.RootCmd.SetArgs([]string{"add", "--category", "Entertainment", "--amount=-50", "--desc", "Games"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Add expense failed: %v", err)
	}

	// 4. Verify Progress Bar in 'budget list'
	listOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"budget", "list"})
	cli.RootCmd.Execute()

	// It should show [█████-----] 50%
	if !strings.Contains(listOut.String(), "50%") {
		t.Errorf("Expected progress bar to show 50%%, got output:\n%s", listOut.String())
	}

	// 5. Overspend (Add 60 more -> Total 110 spent vs 100 limit)
	addOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(addOut)
	cli.RootCmd.SetArgs([]string{"add", "--category", "Entertainment", "--amount=-60", "--desc", "Concert"})
	cli.RootCmd.Execute()

	// 6. Verify Immediate Alert Message in 'add' output
	if !strings.Contains(addOut.String(), "ALERT") {
		t.Errorf("Expected Alert message in add output, got:\n%s", addOut.String())
	}

	// 7. Verify Status in 'budget list'
	listOut.Reset()
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"budget", "list"})
	cli.RootCmd.Execute()

	if !strings.Contains(listOut.String(), "OVER BUDGET") {
		t.Errorf("Expected status 'OVER BUDGET', got output:\n%s", listOut.String())
	}
}
