package cli

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SebiGabor/personal-finance-cli/internal/db"
	"github.com/spf13/cobra"
)

var database *sql.DB

// RootCmd is the base command for the application
var RootCmd = &cobra.Command{
	Use:   "finance",
	Short: "Personal Finance CLI Manager",
	Long:  `A command-line tool for tracking personal income and expenses.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// If the database isn't already set (e.g. by a test), connect to the production DB
		if database == nil {
			var err error
			database, err = db.Connect()
			if err != nil {
				return fmt.Errorf("could not connect to database: %w", err)
			}
		}
		return nil
	},
}

// SetDatabase allows external packages (like tests) to inject a database connection
func SetDatabase(db *sql.DB) {
	database = db
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
