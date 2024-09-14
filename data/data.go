package data

import (
	"buzz/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func GetJobs(db *sql.DB) [][]string {
	var jobs [][]string
	rows, _ := db.Query("SELECT position, company, salary, status FROM jobs")
	defer rows.Close()

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		thisJob := models.Job{}
		err = rows.Scan(&thisJob.Position, &thisJob.Company, &thisJob.Salary, &thisJob.Status)
		if err != nil {
			log.Fatal(err)
		}

		jobs = append(jobs, []string{thisJob.Position, thisJob.Company, thisJob.Salary, thisJob.Status})

	}
	return jobs
}
