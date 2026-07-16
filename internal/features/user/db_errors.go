package user

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func MapDBDuplicateError(err error, errorMap map[string]string) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1062) {
		if strings.Contains(mysqlErr.Message, "unique_username") {
			errorMap["username"] = "Usuario ya registrado."
		}
	}
}