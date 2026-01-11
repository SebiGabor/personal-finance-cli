package cli

import (
	"github.com/SebiGabor/personal-finance-cli/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal UI",
	Long:  "Opens an interactive table to browse and scroll through transactions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.StartTUI(database)
	},
}

func init() {
	RootCmd.AddCommand(tuiCmd)
}
