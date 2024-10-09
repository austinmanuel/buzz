package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

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

func deleteJob() {

}
