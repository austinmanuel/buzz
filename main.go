package main

import (
	"buzz/ui"
	"database/sql"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	db := startDb()
	m := ui.BuildTable(db)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}

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
