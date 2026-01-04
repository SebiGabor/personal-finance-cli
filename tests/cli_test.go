package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/SebiGabor/personal-finance-cli/internal/cli"
)

func TestCLICommands(t *testing.T) {
	// 1. Initialize test database
	db := NewTestDB(t)
	// 2. Inject the test database into the CLI package
	cli.SetDatabase(db)

	// Test 'add' command
	addOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(addOut)
	cli.RootCmd.SetArgs([]string{"add", "--amount", "50.00", "--desc", "Test Expense", "--category", "Testing"})

	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Add command failed: %v", err)
	}

	if !strings.Contains(addOut.String(), "Successfully added transaction") {
		t.Errorf("Expected success message, got: %s", addOut.String())
	}

	// Test 'list' command
	listOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"list"})

	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	output := listOut.String()
	if !strings.Contains(output, "Test Expense") || !strings.Contains(output, "50.00") {
		t.Errorf("Expected list output to contain transaction details, got: %s", output)
	}
}
