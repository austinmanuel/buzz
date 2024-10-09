package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	Table     table.Model
	altscreen bool
	quitting  bool
}

type FormModel struct {
	form *huh.Form
}

type Job struct {
	Id       int
	Position string
	Company  string
	Salary   string
	Status   string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctl-c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Lets work %s!", m.Table.SelectedRow()[1]))
		case " ":
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	return baseStyle.Render(m.Table.View()) + "\n" +
		helpStyle.Render("  space: select • q / ctl-c: quit\n")
}

func NewModel() FormModel {
	return FormModel{
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

func (f FormModel) Init() tea.Cmd {
	return f.form.Init()
}

func (f FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	if x, ok := form.(*huh.Form); ok {
		f.form = x
	}

	return f, cmd
}

func (f FormModel) View() string {
	if f.form.State == huh.StateCompleted {
		class := f.form.GetString("class")
		level := f.form.GetInt("level")
		return fmt.Sprintf("You selected: %s, Lvl. %d", class, level)
	}
	return f.form.View()
}
