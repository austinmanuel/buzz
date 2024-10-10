package main

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"strconv"
)

type tableModel struct {
	table     table.Model
	db        *sql.DB
	altscreen bool
	quitting  bool
	starting  bool
}

type formModel struct {
	form *huh.Form
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
	if m.starting {
		cmd = tea.EnterAltScreen
		m.starting = false
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			id, err := strconv.Atoi(m.table.SelectedRow()[0])
			checkErr(err)
			cursor := m.table.Cursor()
			log.Infof("D key pressed, deleting job %d", id)
			affected := deleteJob(m.db, id)
			log.Infof("Deleted %d records", affected)
			m.table = buildTable(m.db)
			if cursor != 0 {
				m.table.SetCursor(cursor - 1)
			}
			return m, cmd
		case "esc", "q", "ctl-c":
			log.Info("Escape key pressed, quitting app")
			m.quitting = true
			return m, tea.Quit
		case "n":
			log.Info("N key pressed, creating new job")
			createJob(m.db)
			m.table = buildTable(m.db)
			return m, cmd
		case " ", "enter":
			log.Infof("Enter key pressed, would edit job %s", m.table.SelectedRow()[0])
			//if _, err := tea.NewProgram(newModel()).Run(); err != nil {
			//	log.Fatal(err)
			//}
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
	rows := []table.Row{}
	for _, jobRow := range getJobs(m.db) {
		rows = append(rows, jobRow)
	}
	m.table.SetRows(rows)
}

func newModel() formModel {
	return formModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("class").
					Options(huh.NewOptions("Warrior", "Mage", "Rogue")...).
					Title("Choose your class"),

				huh.NewSelect[int]().
					Key("level").
					Options(huh.NewOptions(1, 20, 9999)...).
					Title("Choose your level"),
			),
		),
	}
}

func (f formModel) Init() tea.Cmd {
	return f.form.Init()
}

func (f formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	if x, ok := form.(*huh.Form); ok {
		f.form = x
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
		case "q", "ctl-c":
			return f, tea.Quit
		}
	}
	return f, cmd
}

func (f formModel) View() string {
	if f.form.State == huh.StateCompleted {
		class := f.form.GetString("class")
		level := f.form.GetInt("level")
		return fmt.Sprintf("You selected: %s, Lvl. %d", class, level)
	}
	return f.form.View()
}
