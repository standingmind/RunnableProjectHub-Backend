package dbutil

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Username string
	Password string
	Email    string
}

func CreateUser(con *sql.DB, user User) error {
	query := "insert into user(user_name,user_password,user_email,user_create_date,user_last_modified_date) values(?,?,?,?,?) "
	prpStmt, err := con.Prepare(query)
	if err != nil {
		return err
	}
	defer prpStmt.Close()
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")

	fmt.Println(ts)
	if user.Username != "" && user.Password != "" && user.Email != "" {
		_, err = prpStmt.Exec(user.Username, user.Password, user.Email, ts, ts)
		if err != nil {
			return errors.New("Error Creating User :")
		}
	} else {
		return errors.New("Empty user info")
	}

	fmt.Println("created user successfully")
	return nil
}

func IsEmailExist(email string, con *sql.DB) (bool, error) {
	var ans bool
	query := "Select Exists(Select * from user where user_email=?)"
	prpStmt, err := con.Prepare(query)
	if err != nil {
		return ans, err
	}
	row := prpStmt.QueryRow(email)
	err = row.Scan(&ans)
	if err != nil {
		return ans, err
	}
	return ans, nil
}

func GetUserByEmail(email string, con *sql.DB) (int, string, error) {
	var username = ""
	var userId = -1
	query := "Select user_id,user_name from user where user_email=?"
	prpStmt, err := con.Prepare(query)
	if err != nil {
		return userId, username, err
	}
	row := prpStmt.QueryRow(email)
	err = row.Scan(&userId, &username)
	if err != nil {
		return userId, username, err
	}
	return userId, username, nil
}

func IsUserExist(email string, password string, con *sql.DB) (bool, error) {
	var ans = false
	query := "Select Exists(Select * from user where user_email=? and user_password=?)"
	prpStmt, err := con.Prepare(query)
	if err != nil {
		return ans, err
	}
	row := prpStmt.QueryRow(email, password)
	err = row.Scan(&ans)
	if err != nil {
		return ans, err
	}
	return ans, nil
}
