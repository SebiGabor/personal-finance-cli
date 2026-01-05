package cli

import (
	"fmt"
	"strconv"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a transaction by its ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		err = models.DeleteTransaction(database, id)
		if err != nil {
			return fmt.Errorf("failed to delete transaction: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Transaction %d successfully deleted.\n", id)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
