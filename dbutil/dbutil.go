package dbutil

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBComponent struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Driver   string `json:"driver"`
}

func GetConnection() (*sql.DB, error) {
	dbInfo, err := GetConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dbDriver := dbInfo.Driver
	dbName := dbInfo.Database
	dbUser := dbInfo.Username
	dbPass := dbInfo.Password
	dbHost := dbInfo.Hostname

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@("+dbHost+")/"+dbName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("successfully connected to database")

	return db, nil
}

func CloseConnection(db *sql.DB) error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetConfig() (DBComponent, error) {
	//read component from config.json
	var dbcomponent DBComponent
	configFile, err := os.Open("./config.json")
	if err != nil {
		return dbcomponent, err
	}
	defer configFile.Close()
	readByte, err := ioutil.ReadAll(configFile)
	if err != nil {
		return dbcomponent, err
	}
	json.Unmarshal(readByte, &dbcomponent)

	fmt.Println(dbcomponent.Username)
	fmt.Println(dbcomponent.Database)
	fmt.Println(dbcomponent.Hostname)
	fmt.Println(dbcomponent.Password)
	fmt.Println(dbcomponent.Driver)
	return dbcomponent, nil
}
