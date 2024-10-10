package main

import (
	"database/sql"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

func startDb() *sql.DB {
	if _, err := os.Stat("./data/jobs.db"); os.IsNotExist(err) {
		log.Info("No database file found, creating jobs.db in local directory.")
		createDb()
	}
	db, err := sql.Open("sqlite3", "./data/jobs.db")
	checkErr(err)
	return db
}

func createDb() {
	const create string = "CREATE TABLE IF NOT EXISTS `jobs` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `position` TEXT, `company` TEXT, `salary` TEXT, `status` TEXT)"

	db, err := sql.Open("sqlite3", "./data/jobs.db")
	checkErr(err)

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
	checkErr(err)

	return
}

func getJobs(db *sql.DB) [][]string {
	var jobs [][]string
	rows, _ := db.Query("SELECT id, position, company, salary, status FROM jobs")
	defer rows.Close()

	err := rows.Err()
	checkErr(err)

	for rows.Next() {
		thisJob := job{}
		err = rows.Scan(&thisJob.id, &thisJob.position, &thisJob.company, &thisJob.salary, &thisJob.status)
		checkErr(err)

		jobs = append(jobs, []string{strconv.Itoa(thisJob.id), thisJob.position, thisJob.company, thisJob.salary, thisJob.status})

	}
	return jobs
}

func deleteJob(db *sql.DB, id int) int64 {
	stmt, err := db.Prepare("DELETE FROM jobs WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(id)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}

func createJob(db *sql.DB) int64 {
	stmt, err := db.Prepare("INSERT INTO jobs (id, position, company, salary, status) VALUES (?, ?, ?, ?, ?)")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(nil, "Test Position", "Test Company", "Test Salary", "Test Status")
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}
