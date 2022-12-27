package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/ttanik/bookstore_utils-go/rest_errors"
)

const (
	errorQueryNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorQueryNoRows) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", err)
	}

	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError(fmt.Sprintf("invalid data %s", sqlErr.Message))
	}
	return rest_errors.NewInternalServerError("error processing request", err)
}
