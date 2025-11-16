package tests

import (
	"testing"

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
