package tests

import (
	"testing"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
)

func TestBudgetCRUD(t *testing.T) {
	db := NewTestDB(t)

	b := &models.Budget{
		Category: "Food",
		Amount:   300,
		Period:   "monthly",
	}

	// Create
	if err := models.CreateBudget(db, b); err != nil {
		t.Fatalf("CreateBudget failed: %v", err)
	}

	// Get
	got, err := models.GetBudget(db, b.ID)
	if err != nil {
		t.Fatalf("GetBudget failed: %v", err)
	}
	if got.Amount != 300 {
		t.Errorf("expected amount 300, got %v", got.Amount)
	}

	// Update
	b.Amount = 350
	if err := models.UpdateBudget(db, b); err != nil {
		t.Fatalf("UpdateBudget failed: %v", err)
	}

	got, _ = models.GetBudget(db, b.ID)
	if got.Amount != 350 {
		t.Errorf("expected updated amount 350")
	}

	// Delete
	if err := models.DeleteBudget(db, b.ID); err != nil {
		t.Fatalf("DeleteBudget failed: %v", err)
	}
}

func TestBudgetSpending(t *testing.T) {
	db := NewTestDB(t)

	// 1. Setup Data:
	// - 2 Expenses in "Food" for current month
	// - 1 Income in "Food" (refund? shouldn't count towards spending limit ideally, or should reduce spending.
	//   Current logic ignores positive amounts based on our query 'amount < 0')
	// - 1 Expense in "Food" for PREVIOUS month (should be ignored)

	now := time.Now()

	transactions := []*models.Transaction{
		{Date: now, Description: "Burger", Amount: -10.00, Category: "Food"},
		{Date: now, Description: "Pizza", Amount: -20.00, Category: "Food"},
		{Date: now, Description: "Refund", Amount: 5.00, Category: "Food"},                         // Positive, ignored by current logic
		{Date: now, Description: "Gas", Amount: -50.00, Category: "Transport"},                     // Different category
		{Date: now.AddDate(0, -1, 0), Description: "Old Pizza", Amount: -100.00, Category: "Food"}, // Wrong month
	}

	for _, tr := range transactions {
		if err := models.CreateTransaction(db, tr); err != nil {
			t.Fatalf("failed to setup transaction: %v", err)
		}
	}

	// 2. Test GetSpendingTotal
	total, err := models.GetSpendingTotal(db, "Food", now.Month(), now.Year())
	if err != nil {
		t.Fatalf("GetSpendingTotal failed: %v", err)
	}

	// Expected: 10 + 20 = 30. (5 is income, 50 is transport, 100 is last month)
	if total != 30.00 {
		t.Errorf("expected spending 30.00, got %.2f", total)
	}

	// 3. Test Previous Month (Old Pizza)
	totalLast, _ := models.GetSpendingTotal(db, "Food", now.AddDate(0, -1, 0).Month(), now.AddDate(0, -1, 0).Year())
	if totalLast != 100.00 {
		t.Errorf("expected last month spending 100.00, got %.2f", totalLast)
	}
}
