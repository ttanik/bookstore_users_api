package users

import (
	"fmt"

	"github.com/ttanik/bookstore_users-api/datasources/mysql/users_db"
	"github.com/ttanik/bookstore_users-api/logger"
	"github.com/ttanik/bookstore_users-api/utils/errors"
	"github.com/ttanik/bookstore_users-api/utils/mysql_utils"
)

const (
	errorQueryNoRows      = "no rows in result set"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryInsertUser       = "INSERT into users(first_name, last_name, email, date_created, password, status) VALUES(?,?,?,?,?,?);"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=?, date_created=? WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
		&user.Status,
	); err != nil {
		logger.Error("error when trying to scan user", err)
		return mysql_utils.ParseError(err)
	}
	return nil
}
func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	_, queryErr := stmt.Exec(user.Id)
	if queryErr != nil {
		logger.Error("error when trying to execute delete user", queryErr)
		return mysql_utils.ParseError(queryErr)
	}
	return nil
}
func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.DateCreated,
		user.Password,
		user.Status,
	)
	if err != nil {
		logger.Error("error when trying to execute save user", err)
		return mysql_utils.ParseError(err)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id", err)
		return errors.NewInternalServerError(fmt.Sprintf("error trying to save user %s", err.Error()))
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Id)
	if err != nil {
		logger.Error("error when trying to execute user update", err)
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user statement", err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to execute find user", err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.DateCreated,
			&user.Status,
		); err != nil {
			logger.Error("error when trying to parse user on find user", err)
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
