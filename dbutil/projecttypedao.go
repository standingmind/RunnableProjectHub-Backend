package dbutil

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type ProjectType struct {
	project_type_id                 int
	project_type_name               string
	project_type_create_date        string
	project_type_last_modified_date string
}

func CreateProjectType(con *sql.DB, projectType ProjectType) error {
	query := "Insert into projectType(project_type_name,project_type_create_date,project_type_last_modified_date) into values(?,?,?,?)"
	prpStmt, err := con.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}
	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")
	_, err = prpStmt.Exec(projectType.project_type_name, date, date)
	if err != nil {
		log.Println(err)
		return errors.New("Cannot Create Project Type")
	}

	return nil
}
