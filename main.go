package main

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

type job struct {
	id       int
	position string
	company  string
	salary   string
	status   string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctl-c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Lets work %s!", m.table.SelectedRow()[1]))
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	db := startDb()
	m := buildTable(db)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}

}

func buildTable(db *sql.DB) model {
	columns := []table.Column{
		{Title: "Position", Width: 30},
		{Title: "Company", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Salary", Width: 10},
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

	return model{t}
}

func getJobs(db *sql.DB) [][]string {
	var jobs [][]string
	rows, _ := db.Query("SELECT position, company, salary, status FROM jobs")
	defer rows.Close()

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		thisJob := job{}
		err = rows.Scan(&thisJob.position, &thisJob.company, &thisJob.salary, &thisJob.status)
		if err != nil {
			log.Fatal(err)
		}

		jobs = append(jobs, []string{thisJob.position, thisJob.company, thisJob.salary, thisJob.status})

	}
	return jobs
}

func startDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createDb() *sql.DB {
	fmt.Println("In createDb")
	const create string = `
		CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER, 
		position TEXT, 
		company TEXT, 
		salary TEXT, 
		status TEXT, 
		(PRIMARY KEY id AUTOINCREMENT)
		);`

	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(create); err != nil {
		log.Fatal(err)
	}

	return db
}
