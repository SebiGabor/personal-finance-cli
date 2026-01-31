package tests

import (
	"testing"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
)

func TestNormalizeCategory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Lowercase", input: "food", expected: "Food"},
		{name: "Uppercase", input: "GROCERIES", expected: "Groceries"},
		{name: "Mixed Case", input: "tRanSpoRt", expected: "Transport"},
		{name: "With Spaces", input: "  entertainment  ", expected: "Entertainment"},
		{name: "Empty", input: "", expected: "Uncategorized"},
		{name: "Already Correct", input: "Salary", expected: "Salary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.NormalizeCategory(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeCategory(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
