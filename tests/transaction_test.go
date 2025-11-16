package tests

import (
	"testing"
	"time"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
)

func TestTransactionCRUD(t *testing.T) {
	db := NewTestDB(t)

	// 1. Create
	tr := &models.Transaction{
		Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Test Income",
		Amount:      200.50,
		Category:    "Salary",
	}

	err := models.CreateTransaction(db, tr)
	if err != nil {
		t.Fatalf("CreateTransaction failed: %v", err)
	}
	if tr.ID == 0 {
		t.Fatalf("expected ID to be set")
	}

	// 2. Get
	loaded, err := models.GetTransaction(db, tr.ID)
	if err != nil {
		t.Fatalf("GetTransaction failed: %v", err)
	}
	if loaded.Amount != 200.50 {
		t.Errorf("expected amount 200.50, got %v", loaded.Amount)
	}

	// 3. Update
	tr.Amount = 220.00
	if err := models.UpdateTransaction(db, tr); err != nil {
		t.Fatalf("UpdateTransaction failed: %v", err)
	}

	updated, _ := models.GetTransaction(db, tr.ID)
	if updated.Amount != 220.00 {
		t.Errorf("expected updated amount 220.00")
	}

	// 4. Delete
	if err := models.DeleteTransaction(db, tr.ID); err != nil {
		t.Fatalf("DeleteTransaction failed: %v", err)
	}

	_, err = models.GetTransaction(db, tr.ID)
	if err == nil {
		t.Fatalf("expected error after deleting transaction")
	}
}
