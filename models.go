package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type tableModel struct {
	table     table.Model
	altscreen bool
	quitting  bool
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
	tea.EnterAltScreen()
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			deleteJob()
		case "esc", "q", "ctl-c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			//tea.Printf("Lets work %s!", m.Table.SelectedRow()[1])
			if _, err := tea.NewProgram(newModel()).Run(); err != nil {
				log.Fatal(err)
			}
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
		helpStyle.Render("  enter: select • d: delete • q / ctl-c: quit\n")
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
