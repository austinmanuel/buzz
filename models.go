package main

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"strconv"
)

type tableModel struct {
	table    table.Model
	db       *sql.DB
	quitting bool
}

type job struct {
	id       int
	position string
	company  string
	salary   string
	status   string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			log.Info("D key pressed")
			cursor := m.table.Cursor()
			confirmation := m.deleteJob()
			m.table = buildTable(m.db)
			if (cursor != 0) && (confirmation) {
				m.table.SetCursor(cursor - 1)
			} else {
				m.table.SetCursor(cursor)
			}
			return m, cmd
		case "esc", "q", "ctl-c":
			log.Info("Escape key pressed")
			m.quitting = true
			return m, tea.Quit
		case "n":
			log.Info("N key pressed")
			m.newJob()
			m.table = buildTable(m.db)
			return m, cmd
		case " ", "enter":
			log.Info("Enter key pressed")
			cursor := m.table.Cursor()
			m.updateJob(m.table.SelectedRow())
			m.table = buildTable(m.db)
			m.table.SetCursor(cursor)
			return m, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	return baseStyle.Render(m.table.View()) + "\n" +
		helpStyle.Render("  enter / space: select • n: new • d: delete • q / ctl-c: quit\n")
}

func (m tableModel) updateRows() {
	rows := make([]table.Row, 0)
	for _, jobRow := range getJobEntries(m.db) {
		rows = append(rows, jobRow)
	}
	m.table.SetRows(rows)
}

func (m tableModel) newJob() {
	input := jobForm()
	err := input.Run()
	checkErr(err)

	newJob := job{
		0,
		input.GetString("position"),
		input.GetString("company"),
		input.GetString("salary"),
		input.GetString("status"),
	}
	log.Infof("New Job - Position: %s, Company: %s, Salary: %s, Status: %s",
		newJob.position, newJob.company, newJob.salary, newJob.status)

	createJobEntry(m.db, newJob)
}

func (m tableModel) deleteJob() bool {
	confirmForm := confirm()
	err := confirmForm.Run()
	checkErr(err)

	if confirmForm.GetBool("confirmation") {
		id, err := strconv.Atoi(m.table.SelectedRow()[0])
		checkErr(err)
		affected := deleteJobEntry(m.db, id)
		log.Infof("Deleted %d records", affected)
		return true
	}
	return false
}

func (m tableModel) updateJob(oldJobData table.Row) {
	id, err := strconv.Atoi(oldJobData[0])
	checkErr(err)
	oldJob := job{
		id,
		oldJobData[1],
		oldJobData[2],
		oldJobData[3],
		oldJobData[4],
	}
	log.Infof("Updating Job - Position: %s, Company: %s, Salary: %s, Status: %s",
		oldJob.position, oldJob.company, oldJob.salary, oldJob.status)

	update := updateForm(oldJob)
	err = update.Run()
	checkErr(err)
	updatedJob := job{
		0,
		update.GetString("position"),
		update.GetString("company"),
		update.GetString("salary"),
		update.GetString("status"),
	}
	affected := updateJobEntry(m.db, updatedJob, id)
	log.Infof("Updated %d records", affected)
}
