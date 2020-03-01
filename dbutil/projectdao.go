package dbutil

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Project struct {
	ID                 int
	Title              string
	Language           string
	LanguageVersion    string
	Server             string
	ServerVersion      string
	Database           string
	DatabaseVersion    string
	DatabaseBackupPath string
	Description        string
	Path               string
	CreateDate         string
	LastModifiedDate   string
	UserID             int
	ProjectType        string
	ZipPath            string
	Permission         string
	Device             string
	AndroidVersion     string
	APIVersion         string
	ContainerName      string
	ContainerIP        string
	DatabaseName       string
}

type ProjectInfos struct {
	Infos    []Info `json:"projectInfos"`
	UserID   string `json:"userId"`
	Username string `json:"userName"`
	Email    string `json:"email"`
	// ContainerExist func(name string) bool
	// ContainerStart func(name string)
	// ContainerIP    func(name string) string
	// ContainerStop  func(name string)

}

type Info struct {
	ProjectName        string `json:"projectName"`
	ContainerName      string `json:"containerName"`
	ProjectDescription string `json:"projectDescription"`
	IPAddress          string `json:"IPAddress"`
	ProjectType        string `json:"projectType"`
	ZipName            string `json:"zipName"`
	ProjectID          string `json:"projectId"`
	UserID             string `json:"usreId"`
	Permission         string `json:"permission"`
}

func GetLastId(con *sql.DB) (int, error) {
	var lastId int
	query := "Select AUTO_INCREMENT from information_schema.TABLES where TABLE_SCHEMA=? and TABLE_NAME=?"

	prpStmt, err := con.Prepare(query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer prpStmt.Close()
	row := prpStmt.QueryRow("runnableprojecthub", "project")
	err = row.Scan(&lastId)
	if err != nil {
		log.Println(err)
		return lastId, err
	}
	return lastId, nil

}

func SaveWebAndWindowProject(con *sql.DB, project Project) error {
	query := `Insert into project(project_title,
		project_language,
		project_language_version,
		project_database,
		project_database_version,
		project_database_backup_path,
		project_description,
		project_path,
		project_create_date,
		project_last_modified_date,
		user_id,
		project_permission,
		project_type,
		project_zip_path,
		project_container_name,
		project_container_ip,
		project_database_name) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	prpStmt, err := con.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer prpStmt.Close()

	_, err = prpStmt.Exec(project.Title,
		project.Language,
		project.LanguageVersion,
		project.Database,
		project.DatabaseVersion,
		project.DatabaseBackupPath,
		project.Description,
		project.Path,
		project.CreateDate,
		project.LastModifiedDate,
		project.UserID,
		project.Permission,
		project.ProjectType,
		project.ZipPath,
		project.ContainerName,
		project.ContainerIP,
		project.DatabaseName)

	if err != nil {
		log.Println(err)
		return errors.New("cannot create project")
	}

	fmt.Println("user created successfully ! ")
	return nil
}

func SaveAndroidProject(con *sql.DB, project Project) error {
	query := `Insert into project(project_title,
		project_description,
		project_path,
		project_create_date,
		project_last_modified_date,
		user_id,
		project_permission,
		project_type,
		project_zip_path,
		project_device,
		project_android_version,
		project_api_version,
		project_container_name,
		project_container_ip,
		project_database_name) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	prpStmt, err := con.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer prpStmt.Close()

	_, err = prpStmt.Exec(project.Title,
		project.Description,
		project.Path,
		project.CreateDate,
		project.LastModifiedDate,
		project.UserID,
		project.Permission,
		project.ProjectType,
		project.ZipPath,
		project.Device,
		project.AndroidVersion,
		project.APIVersion,
		project.ContainerName,
		project.ContainerIP,
		project.DatabaseName)

	if err != nil {
		log.Println(err)
		return errors.New("cannot create project")
	}

	fmt.Println("user created successfully ! ")
	return nil
}

func GetProjectInfo(con *sql.DB, pattern string, projectType string, language string) ([]Info, error) {
	infoArray := []Info{}
	var rows *sql.Rows
	var query string
	if language == "All" && projectType == "All" {
		query = "Select project_title,project_container_name,project_description,project_container_ip,project_type,project_zip_path,project_id,user_id,project_permission from project where project_title like '%" + pattern + "%'"
	} else if projectType == "All" {
		query = "Select project_title,project_container_name,project_description,project_container_ip,project_type,project_zip_path,project_id,user_id,project_permission from project where project_title like '%" + pattern + "%' and project_language = '" + language + "'"
	} else if language == "All" {
		query = "Select project_title,project_container_name,project_description,project_container_ip,project_type,project_zip_path,project_id,user_id,project_permission from project where project_title like '%" + pattern + "%' and project_type = '" + projectType + "'"
	} else {
		query = "Select project_title,project_container_name,project_description,project_container_ip,project_type,project_zip_path,project_id,user_id,project_permission from project where project_title like '%" + pattern + "%' and project_type ='" + projectType + "' or project_language = '" + language + "'"
	}

	fmt.Println(query)

	prpStmt, err := con.Prepare(query)
	if err != nil {
		log.Println(err)
		return infoArray, err
	}
	defer prpStmt.Close()
	rows, _ = prpStmt.Query()

	// Fetch rows
	for rows.Next() {
		var row Info
		err = rows.Scan(&row.ProjectName, &row.ContainerName, &row.ProjectDescription, &row.IPAddress, &row.ProjectType, &row.ZipName, &row.ProjectID, &row.UserID, &row.Permission)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		infoArray = append(infoArray, row)
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(infoArray)
	return infoArray, nil
}
