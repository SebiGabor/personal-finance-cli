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

func TestSearchAndDelete(t *testing.T) {
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 1. Add a transaction to search for
	cli.RootCmd.SetArgs([]string{"add", "--amount", "15.00", "--desc", "Pizza", "--category", "Food"})
	cli.RootCmd.Execute()

	// 2. Test Search
	searchOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(searchOut)
	cli.RootCmd.SetArgs([]string{"search", "Pizza"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if !strings.Contains(searchOut.String(), "Pizza") {
		t.Errorf("Search should find 'Pizza', got: %s", searchOut.String())
	}

	// 3. Test Delete (ID should be 1 as it's a fresh DB)
	delOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(delOut)
	cli.RootCmd.SetArgs([]string{"delete", "1"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !strings.Contains(delOut.String(), "successfully deleted") {
		t.Errorf("Delete success message missing, got: %s", delOut.String())
	}
}
