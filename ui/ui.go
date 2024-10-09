package ui

import (
	"buzz/data"
	"buzz/models"
	"database/sql"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func BuildTable(db *sql.DB) models.TableModel {
	columns := []table.Column{
		{Title: "Position", Width: 30},
		{Title: "Company", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Salary", Width: 10},
	}

	rows := []table.Row{}

	for _, jobRow := range data.GetJobs(db) {
		rows = append(rows, jobRow)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return models.TableModel{t}
}
