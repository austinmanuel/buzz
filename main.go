package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func main() {
	db := startDb()
	m := buildTable(db)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}

}
