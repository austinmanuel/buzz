package main

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func buildTableModel(db *sql.DB) tableModel {
	return tableModel{buildTable(db), db, true, false, true}
}

func buildTable(db *sql.DB) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Position", Width: 30},
		{Title: "Company", Width: 20},
		{Title: "Salary", Width: 10},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}

	for _, jobRow := range getJobs(db) {
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

	return t
}
