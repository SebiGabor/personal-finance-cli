package tui

import (
	"database/sql"
	"fmt"

	"github.com/SebiGabor/personal-finance-cli/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StartTUI launches the interactive terminal interface
func StartTUI(db *sql.DB) error {
	app := tview.NewApplication()

	// 1. Fetch Data
	transactions, err := models.ListTransactions(db)
	if err != nil {
		return err
	}

	// 2. Create Table
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false). // Select rows, not individual cells
		SetFixed(1, 0)              // Fix the header row

	// 3. Set Headers
	headers := []string{"ID", "DATE", "CATEGORY", "DESCRIPTION", "AMOUNT"}
	for i, h := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	// 4. Populate Rows
	for i, t := range transactions {
		row := i + 1

		// Color logic: Green for Income, Red for Expense
		color := tcell.ColorWhite
		if t.Amount < 0 {
			color = tcell.ColorRed
		} else {
			color = tcell.ColorGreen
		}

		// ID
		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", t.ID)).SetAlign(tview.AlignCenter))

		// Date
		table.SetCell(row, 1, tview.NewTableCell(t.Date.Format("2006-01-02")).SetAlign(tview.AlignCenter))

		// Category
		table.SetCell(row, 2, tview.NewTableCell(t.Category).SetAlign(tview.AlignCenter))

		// Description (Limit length to keep UI clean)
		desc := t.Description
		if len(desc) > 30 {
			desc = desc[:27] + "..."
		}
		table.SetCell(row, 3, tview.NewTableCell(desc))

		// Amount
		table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
	}

	// 5. Layout & Keybindings
	// Title
	frame := tview.NewFrame(table).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Personal Finance Manager", true, tview.AlignCenter, tcell.ColorGreen).
		AddText("Press 'q' or 'Esc' to quit | Use Arrow Keys to navigate", false, tview.AlignCenter, tcell.ColorGray)

	// Quit functionality
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(frame, true).SetFocus(table).Run(); err != nil {
		return err
	}

	return nil
}
