package main

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func buildTableModel(db *sql.DB) tableModel {
	return tableModel{buildTable(db), db, false}
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

	for _, jobRow := range getJobEntries(db) {
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

func jobForm() *huh.Form {
	var (
		position string
		company  string
		salary   string
		status   string
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("position").Prompt("Job Title: ").Value(&position),
			huh.NewInput().Key("company").Prompt("Company: ").Value(&company),
			huh.NewInput().Key("salary").Prompt("Salary: ").Value(&salary),
			huh.NewInput().Key("status").Prompt("Status: ").Value(&status),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("Finished?").Affirmative("Yes").Negative("No"),
		),
	)
	return form
}

func updateForm(oldJob job) *huh.Form {
	var (
		position string
		company  string
		salary   string
		status   string
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("position").Prompt("Job Title: ").Placeholder(oldJob.position).Value(&position),
			huh.NewInput().Key("company").Prompt("Company: ").Placeholder(oldJob.company).Value(&company),
			huh.NewInput().Key("salary").Prompt("Salary: ").Placeholder(oldJob.salary).Value(&salary),
			huh.NewInput().Key("status").Prompt("Status: ").Placeholder(oldJob.status).Value(&status),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("Finished?").Affirmative("Yes").Negative("No"),
		),
	)
	return form
}

func confirm() *huh.Form {
	var confirmation bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to delete?").
				Affirmative("Yes").
				Negative("No").Value(&confirmation),
		),
	)
	return form
}
