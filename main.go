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
	if _, err := os.Stat("/data/jobs.db"); os.IsNotExist(err) {
		createDb()
	}
	db, err := sql.Open("sqlite3", "./data/jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createDb() {
	const create string = "CREATE TABLE IF NOT EXISTS `jobs` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `position` TEXT, `company` TEXT, `salary` TEXT, `status` TEXT)"

	db, err := sql.Open("sqlite3", "./data/jobs.db")
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if _, err := db.Exec(create); err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

	return
}
