package tests

import (
	"bytes"
	"os"
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

func TestImportCommand(t *testing.T) {
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 1. Create a temporary CSV file
	content := []byte("date,description,amount,category\n2024-05-01,Gym,-30.00,Health\n2024-05-02,Bonus,500.00,Income")
	tmpfile, err := os.CreateTemp("", "transactions-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 2. Run Import Command
	importOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(importOut)
	cli.RootCmd.SetArgs([]string{"import", tmpfile.Name()})

	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Import command failed: %v", err)
	}

	if !strings.Contains(importOut.String(), "2 imported") {
		t.Errorf("Expected '2 imported', got: %s", importOut.String())
	}

	// 3. Verify they are in the DB via List
	listOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"list"})
	cli.RootCmd.Execute()

	if !strings.Contains(listOut.String(), "Gym") || !strings.Contains(listOut.String(), "Bonus") {
		t.Errorf("List output missing imported items")
	}
}

func TestReportCommand(t *testing.T) {
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 1. Add some transactions for a specific month (e.g., May 2024)
	// Income
	cli.RootCmd.SetArgs([]string{"add", "--amount", "1000", "--desc", "Salary", "--category", "Income", "--date", "2024-05-01"})
	cli.RootCmd.Execute()
	// Expense
	cli.RootCmd.SetArgs([]string{"add", "--amount", "-200", "--desc", "Groceries", "--category", "Food", "--date", "2024-05-05"})
	cli.RootCmd.Execute()

	// 2. Run Report for that month
	reportOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(reportOut)
	cli.RootCmd.SetArgs([]string{"report", "--year", "2024", "--month", "5"})

	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Report command failed: %v", err)
	}

	output := reportOut.String()

	// 3. Verify Output
	if !strings.Contains(output, "Total Income:      1000.00") {
		t.Errorf("Expected Total Income 1000.00, got output:\n%s", output)
	}
	if !strings.Contains(output, "Food") {
		t.Errorf("Expected category 'Food' in report")
	}
}

func TestAutoCategorization(t *testing.T) {
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 1. Add a Rule: "Uber" -> "Transport"
	cli.RootCmd.SetArgs([]string{"rules", "add", "--pattern", "(?i)uber", "--category", "Transport"})
	cli.RootCmd.Execute()

	// 2. Create a CSV with an uncategorized "Uber" transaction
	content := []byte("date,description,amount\n2024-06-01,Uber Trip,-25.50")
	tmpfile, err := os.CreateTemp("", "autocat-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(content)
	tmpfile.Close()

	// 3. Import
	cli.RootCmd.SetArgs([]string{"import", tmpfile.Name()})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// 4. Verify the category was set to "Transport"
	listOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"list"})
	cli.RootCmd.Execute()

	if !strings.Contains(listOut.String(), "Transport") {
		t.Errorf("Auto-categorization failed. Expected 'Transport', got:\n%s", listOut.String())
	}
}

func TestBudgetCommands(t *testing.T) {
	db := NewTestDB(t)
	cli.SetDatabase(db)

	// 1. Set a Budget
	cli.RootCmd.SetArgs([]string{"budget", "add", "--category", "Groceries", "--amount", "500"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Budget add failed: %v", err)
	}

	// 2. List Budgets and verify output
	listOut := new(bytes.Buffer)
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"budget", "list"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Budget list failed: %v", err)
	}

	output := listOut.String()
	if !strings.Contains(output, "Groceries") || !strings.Contains(output, "500.00") {
		t.Errorf("Expected budget list to contain 'Groceries' and '500.00', got:\n%s", output)
	}

	// 3. Remove the Budget (ID should be 1)
	cli.RootCmd.SetArgs([]string{"budget", "remove", "1"})
	if err := cli.RootCmd.Execute(); err != nil {
		t.Fatalf("Budget remove failed: %v", err)
	}

	// 4. Verify list is empty
	listOut.Reset()
	cli.RootCmd.SetOut(listOut)
	cli.RootCmd.SetArgs([]string{"budget", "list"})
	cli.RootCmd.Execute()

	if !strings.Contains(listOut.String(), "No budgets set") {
		t.Errorf("Expected 'No budgets set', got:\n%s", listOut.String())
	}
}
