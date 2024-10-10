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

func getJobEntries(db *sql.DB) [][]string {
	var jobs [][]string

	rows, _ := db.Query("SELECT id, position, company, salary, status FROM jobs")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		checkErr(err)
	}(rows)

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

func deleteJobEntry(db *sql.DB, id int) int64 {
	stmt, err := db.Prepare("DELETE FROM jobs WHERE id = ?")
	checkErr(err)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		checkErr(err)
	}(stmt)

	res, err := stmt.Exec(id)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}

func createJobEntry(db *sql.DB, j job) int64 {
	stmt, err := db.Prepare("INSERT INTO jobs (id, position, company, salary, status) VALUES (?, ?, ?, ?, ?)")
	checkErr(err)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		checkErr(err)
	}(stmt)

	res, err := stmt.Exec(nil, j.position, j.company, j.salary, j.status)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}

func updateJobEntry(db *sql.DB, j job, id int) int64 {
	stmt, err := db.Prepare("UPDATE jobs set position = ?, company = ?, salary = ?, status = ? WHERE id = ?")
	checkErr(err)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		checkErr(err)
	}(stmt)

	res, err := stmt.Exec(j.position, j.company, j.salary, j.status, id)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}
