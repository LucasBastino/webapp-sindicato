package validation

import (
	"slices"
	"strings"
	"time"
)

func ValidateNameOrLastName(input string) string {
	// no pongo trim porque borro los espacios despues

    // borro los espacios y despues verifico que no este vacio
	input = removeChars(input, " '")

	if input == "" { 
		return "Campo requerido."
	}	

	if errMsg := isAlphabetic(input); errMsg != "" {
		return errMsg
	}
	
	return isNotLongerThan(input, 50)
}


func ValidateDni(input string) string {
	if input == "" { 
		return ""
	}
	
	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}
	return isNotLongerThan(input, 8)
}

func ValidateBirthday(input string) string {
	if input == "" { 
		return "Campo requerido."
	}
	min := time.Date(1900, 0, 0, 0, 0, 0, 0, time.UTC)
	max := time.Now()
	return isValidDate(input, min, max)
}

func ValidateGender(input string) string {
	if input == "" { 
		return "Campo requerido."
	}

	optionsSet := map[string]struct{}{
		"Masculino":        {},
		"Femenino":         {},
		"No binario":       {},
		"Masculino transgénero":  {},
		"Femenino transgénero":   {},
		"Otra identidad":   {},
		"Sin especificar": {},
	}

	return isValidOption(input, optionsSet)
}

func ValidateRelationship(input string) string {
	if input == "" { 
		return "Campo requerido."
	}

	if errMsg := isAlphabetic(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 20)
}

func ValidateMaritalStatus(input string) string {
	if input == "" { 
		return "Campo requerido."
	}

	optionsSet := map[string]struct{}{
		"Soltero": {},
		"Casado": {},
		"Separado": {},
		"Divorciado": {},
		"Viudo": {},
	}
	return isValidOption(input, optionsSet)
}


// ADDRESS INFO

func ValidateAddress(input string) string {
	if input == "" {
		return ""
	}

	input = removeChars(input, " #-'&~,.º")

	if errMsg := isAlphanumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 100)
}

func ValidateDistrict(input string) string {
	if input == "" {
		return ""
	}

	input = removeChars(input, " #-'&~,.º")

	if errMsg := isAlphanumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 50)
}

func ValidatePostalCode(input string) string {
	if input == "" {
		return ""
	}
	if len(input) == 4{
		if errMsg := isNumeric(input); errMsg != "" {
			return errMsg
		}
		return ""
	}
	
	runes := []rune(input)
	if len(runes) == 8{	
		if errMsg := isAlphabetic(string(runes[0])); errMsg != "" {
			return errMsg
		}
		
		for i := 1; i <= 4; i++ {
			if errMsg := isNumeric(string(runes[i])); errMsg != ""{
				return errMsg
			}
		}
		
		for i := 5; i <= 7; i++ {
			if errMsg := isAlphabetic(string(runes[i])); errMsg != ""{
				return errMsg
			}
		}
		return ""
	}
	return "Debe contener 4 u 8 caracteres."
}


// CONTACT INFO

func ValidatePhone(input string) string {
	if input == "" {
		return ""
	}

	input = removeChars(input, " ")

	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}
	return isNotLongerThan(input, 20)
}

// todo: cambiar en js
func ValidateEmail(input string) string {
	if input == "" {
		return ""
	}


	parts := strings.Split(input, "@")
	if len(parts) != 2{
		return "Formato no válido."
	}
	local, domain := parts[0], parts[1]
	if (len(local) == 0 ||
		!strings.Contains(input, ".") ||
		len(domain) < 3 ||
		// tener en cuenta, "." asi es un string y no se puede comparar
		// '.' asi es un rune / byte literal y sí se puede
		domain[0] == '.' ||
		domain[len(domain)-1] == '.') {
			return "Formato no válido."
		}
	
	domainParts := strings.Split(domain, ".")
	if slices.Contains(domainParts, ""){
		return "Formato no válido." // cubre "".gmail.com", "gmail."", "gmail..com"
	}
	// si no tiene algo antes y despues del punto ||
	// si la ultima parte (el TLD) no tiene minimo 2 caracteres no es válido
	if len(domainParts) < 2 || len(domainParts[len(domainParts)-1]) < 2 {
		return "Formato no válido."
	}
	domain = removeChars(domain, ".")
	if errMsg := isAlphanumeric(domain); errMsg != "" {
		return errMsg
	}

	input = removeChars(input, "!#$%&'*+-/=?^_`{|}~.@")
	if errMsg := isAlphanumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 50)
}

// SOCIAL SECURITY INFO

func ValidateMemberNumber(input string) string {
	if input == "" {
		return ""
	}

	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 20)
}

func ValidateCompanyID(input string) string {
	if input == "" {
		return "Campo requerido."
	}
	
	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}
	if input == "0" {
		return "El valor ingresado no es válido."
	}
	return isNotLongerThan(input, 10)
}



func ValidateCategory(input string) string {
	if input == ""{
		return "Campo requerido."
	}

	optionsSet := map[string]struct{}{
		"Nivel 1: Oficial Múltiple": {},
		"Nivel 2: Oficial Especializado": {},
		"Nivel 3: Oficial General": {},
		"Nivel 4: Medio Oficial": {},
		"Nivel 5: Ayudante": {},
		"Nivel 6: Operario Act. Industrial": {},
	}
	return isValidOption(input, optionsSet)
}

func ValidateEntryDate(input string) string {
	if input == "" {
		return ""
	}
	min := time.Date(1900, 0, 0, 0, 0, 0, 0, time.UTC)
	max := time.Now()
	return isValidDate(input, min, max)
}

func ValidateCompanyName(input string) string {
	if input == "" {
		return "Campo requerido."
	}
	return isNotLongerThan(input, 150)
}

func ValidateContact(input string) string {
	if input == "" {
		return ""
	}
	return isNotLongerThan(input, 200)
}


func ValidateCompanyNumber(input string) string {
	if input == "" {
		return ""
	}

	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 20)
}


func ValidateCuilCuit(input string) string {
	if input == "" { 
		return ""
	}

	input = removeChars(input, " ")

	if strings.Contains(input, "-"){
		parts := strings.Split(input, "-")
		if len(parts) != 3{
			return "Formato no válido."
		}
		prefix, dni, digit := parts[0], parts[1], parts[2]
		if (len(prefix) != 2 ||	len(dni) != 8 || len(digit) != 1){
			return "Formato no válido."
		} else{
			return ""
		}

	}
	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}
	return hasXLength(input, 11)
}


// func ValidateMonth(input string) string {
// 	if input == "" { 
// 		return "Campo requerido."
// 	}

// 	if errMsg := isNumeric(input); errMsg != "" {
// 		return errMsg
// 	}

// 	inputInt, err := strconv.Atoi(input)
// 	if err != nil {
// 		return "El valor ingresado no es válido."
// 	}

// 	if inputInt < 1 || inputInt > 12 {
// 		return "El valor ingresado no es válido."
// 	}

// 	return ""
// }

// func ValidateYear(input string) string {
// 	if input == "" { 
// 		return "Campo requerido."
// 	}

// 	if errMsg := isNumeric(input); errMsg != "" {
// 		return errMsg
// 	}

// 	inputInt, err := strconv.Atoi(input)
// 	if err != nil {
// 		return "El valor ingresado no es válido."
// 	}

// 	if inputInt > time.Now().Year() || inputInt < 2000 {
// 		return "El valor ingresado no es válido."
// 	}

// 	return ""
// }

// func ValidateIsPaid(input string) string {
// 	if input == ""{
// 		return "Campo requerido."
// 	}
// 	if input != "true" && input != "false" {
// 		return "El valor ingresado no es válido."
// 	}
// 	return ""
// }

func ValidateStatus(input string) string {
	if input == ""{
		return "Campo requerido."
	}

	optionsSet := map[string]struct{}{
		"Completado": {},
		"Pendiente": {},
		"Vencido": {},
		"En plan de pago": {},
	}
	return isValidOption(input, optionsSet)
}


func ValidateAmount(input string) string {
	if input == "" {
		return ""
	}

	normalized := NormalizeAmountInput(input)

	if strings.Contains(normalized, ".") {
		if strings.HasPrefix(normalized, ".") || strings.HasSuffix(normalized, ".") {
			return "Formato no válido."
		}

		parts := strings.Split(normalized, ".")
		if len(parts) != 2 {
			return "Formato no válido."
		}

		// no acepta mas de 2 decimales
		if len(parts[1]) > 2 {
			return "Formato no válido."
		}
	}

	digits := removeChars(normalized, ".")

	if errMsg := isNumeric(digits); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(digits, 20)
}

// NormalizeAmountInput converts a user amount string to a ParseFloat-ready value.
// With comma: Argentine style (dots = thousands, comma = decimal).
// Without comma: dot is treated as decimal separator.
func NormalizeAmountInput(input string) string {
	input = strings.TrimSpace(input)
	if strings.Contains(input, ",") {
		input = strings.ReplaceAll(input, ".", "")
		return strings.ReplaceAll(input, ",", ".")
	}
	return input
}



func ValidatePaidAt(input string) string {
	if input == "" {
		return ""
	}
	min := time.Date(1900, 0, 0, 0, 0, 0, 0, time.UTC)
	max := time.Now()
	return isValidDate(input, min, max)
}

func ValidateNumberOfInstallments(input string) string {
	if input == ""{
		return "Campo requerido."
	}

	if errMsg := isNumeric(input); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 2000)
}

func ValidateFirstDueDate(input string) string {
	if input == ""{
		return "Campo requerido."
	}

	min := time.Date(1900, 0, 0, 0, 0, 0, 0, time.UTC)
	max := time.Now().AddDate(10, 0, 0)
	return isValidDate(input, min, max)
}

func ValidateObservations(input string) string {
	return isNotLongerThan(input, 2000)
}



func ValidateUsername(input string) string{
	if input == "" { 
		return "Campo requerido."
	}
	if errMsg := isAlphanumeric(input); errMsg != "" {
		return errMsg
	}

	if errMsg := hasAtLeast(input, 3); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 20)
}

func ValidatePassword(input string) string{
	if input == "" { 
		return "Campo requerido."
	}
	input = removeChars(input, " #-'&,.!?*+")

	if errMsg := isAlphanumeric(input); errMsg != "" {
		return errMsg
	}

	if errMsg := hasAtLeast(input, 8); errMsg != "" {
		return errMsg
	}

	return isNotLongerThan(input, 20)
}

func ValidateAdmin(input string) string {
	if input != "on" && input != ""{
		return "El valor ingresado no es válido."
	}
	return ""
}


func ValidatePermissions(input string) string {
	return isValidOption(input, permissionOptionSet)
}

var permissionOptionSet = map[string]struct{}{
		"editor": {},
		"viewer": {},
		"": {},
	}


/* func FormatToYYYYMMDD(date string) string {
	day := date[0:2]
	month := date[3:5]
	year := date[6:]
	date = year + "/" + month + "/" + day
	return date
}

func FormatToDDMMYYYY(date string) string {
	year := date[0:4]
	month := date[5:7]
	day := date[8:]
	date = day + "/" + month + "/" + year
	return date
} */



// func ValidateCompanyNumber(input, idCompany int) string {

// aca yo despues chequeo que el idcompany sea el mismo y esta todo bien, sino false
// 	if input == "" {
// 		customError.ValidationError.Msg = 'company-number' can't be empty"
// 		return customError.ValidationError
// 	}
// 	if oldCompanyNumber == input {
// 		return nil
// 	}
// 	companiesNumbers, err := GetAllCompaniesNumbersFromDB()
// 	if err != nil {
// 		return err
// 	}
// 	if slices.Contains(companiesNumbers, input) {
// 		customError.ValidationError.Msg = "company number already exists"
// 		return customError.ValidationError
// 	}

// 	return isNumeric(input, "company-number", "")
// }


/* func ValidateAffiliated(input string) string {
	if input == ""{return "El campo no puede estar vacío."}
	inputBool, err := strconv.ParseBool(input)
	if err != nil {
		return "El campo contiene un dato inválido."
	}
	if inputBool || !inputBool {
		return ""
	} else {
		return "El campo contiene un dato inválido."
	}
} */