package tests

import (
	"testing"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
)

func TestCategoryRuleCRUD(t *testing.T) {
	db := NewTestDB(t)

	rule := &models.CategoryRule{
		Pattern:  "(?i)netflix",
		Category: "Entertainment",
	}

	// Create
	if err := models.CreateRule(db, rule); err != nil {
		t.Fatalf("CreateRule failed: %v", err)
	}

	// List
	rules, err := models.ListRules(db)
	if err != nil {
		t.Fatalf("ListRules failed: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	// Delete
	if err := models.DeleteRule(db, rule.ID); err != nil {
		t.Fatalf("DeleteRule failed: %v", err)
	}
}
