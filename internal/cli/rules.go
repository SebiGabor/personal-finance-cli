package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	rulePattern  string
	ruleCategory string
)

// rulesCmd represents the base command for rule management
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage auto-categorization rules",
}

var rulesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new categorization rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		normalizedCategory := models.NormalizeCategory(ruleCategory)
		rule := &models.CategoryRule{
			Pattern:  rulePattern,
			Category: normalizedCategory,
		}

		if err := models.CreateRule(database, rule); err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Rule added: matches '%s' -> '%s'\n", rulePattern, normalizedCategory)
		return nil
	},
}

var rulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categorization rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := models.ListRules(database)
		if err != nil {
			return fmt.Errorf("failed to list rules: %w", err)
		}

		if len(rules) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No rules found.")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tPATTERN\tCATEGORY")
		for _, r := range rules {
			fmt.Fprintf(w, "%d\t%s\t%s\n", r.ID, r.Pattern, r.Category)
		}
		return w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(rulesCmd)
	rulesCmd.AddCommand(rulesAddCmd)
	rulesCmd.AddCommand(rulesListCmd)

	// Flags for 'rules add'
	rulesAddCmd.Flags().StringVarP(&rulePattern, "pattern", "p", "", "Regex pattern to match (e.g., '(?i)netflix')")
	rulesAddCmd.Flags().StringVarP(&ruleCategory, "category", "c", "", "Category to assign")
	rulesAddCmd.MarkFlagRequired("pattern")
	rulesAddCmd.MarkFlagRequired("category")
}
