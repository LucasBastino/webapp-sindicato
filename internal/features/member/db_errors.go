package member

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func mapDBDuplicateError(err error, errorMap map[string]string) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1062) {
		if strings.Contains(mysqlErr.Message, "unique_dni_gender") {
			errorMap["dni"] = "Número de DNI ya registrado."
		}
		if strings.Contains(mysqlErr.Message, "unique_member_number") {
			errorMap["memberNumber"] = "Número de afiliado ya registrado."
		}
		if strings.Contains(mysqlErr.Message, "unique_cuil") {
			errorMap["cuil"] = "Número de CUIL ya registrado."
		}
	}
}