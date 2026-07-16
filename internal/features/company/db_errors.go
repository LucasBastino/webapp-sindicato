package company

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func mapDBDuplicateError(err error, errorMap map[string]string) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1062) {
		if strings.Contains(mysqlErr.Message, "unique_company_number") {
			errorMap["companyNumber"] = "Número de empresa ya registrado."
		}
		if strings.Contains(mysqlErr.Message, "unique_cuit") {
			errorMap["cuit"] = "Número de CUIT ya registrado."
		}
	}
}