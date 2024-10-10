package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func main() {
	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	checkErr(err)
	log.SetOutput(f)

	defer func(f *os.File) {
		err := f.Close()
		checkErr(err)
	}(f)

	db := startDb()
	m := buildTableModel(db)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error encountered: %v", err)
	}
}
