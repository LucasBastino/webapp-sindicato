package creators

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/features/member"
	"github.com/jmoiron/sqlx"
)

// func CreateMembers(c *fiber.Ctx) error {
// 	var CompanyID int
// 	for i := 1; i < 400; i++ {
// 		Name := fmt.Sprintf("nombre%d", i)
// 		// LastName := fmt.Sprintf("apellido%d", i)
// 		DNI := rand.IntN(30000000) + 20000000
// 		// Birthday :=
// 		// Gender := casado soltero y demas
// 		// Phone := fmt.Sprintf("15%d", rand.IntN(99999999))
// 		// hacer nombre+apellido+numero@gmail.com
// 		// Email := fmt.Sprintf("email%d", i)
// 		// Address := fmt.Sprintf("direccion%d", rand.IntN(9999))
// 		// PostalCode := strconv.Itoa(rand.IntN(9999))
// 		// District := fmt.Sprintf("Distrito%d", i)
// 		// MemberNumber := strconv.Itoa(rand.IntN(99999999))
// 		CompanyID = rand.IntN(100) + 1
// 		insert, err := database.DB.Query(fmt.Sprintf("INSERT INTO MemberTable (Name, DNI, CompanyID) VALUES ('%s','%d', '%d')", Name, DNI, CompanyID))
// 		if err != nil {
// 			// DBError{"INSERT MEMBER"}.Error(err)
// 			fmt.Println("error insertando en la DB")
// 			panic(err)
// 		}
// 		insert.Close()
// 	}
// 	return c.SendString("hecho")
// }

func CreateMembers(db *sqlx.DB) error {
	type StreetData struct {
		Code string
		Name string
	}
	type JSONData struct {
		FemaleFirstNames []string
		MaleFirstNames   []string
		LastNames        []string
		MaritalStatus    []string
		Genders          []string
		Streets          []StreetData
	}

	file, err := os.Open("./internal/bootstrap/creators/data.json")
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}
	decoder := json.NewDecoder(file)
	jsonData := JSONData{}
	decoder.Decode(&jsonData)

	members := make([]member.Member, 0, 100)
	for range 100 {
		var m member.Member
		// para que sea mas probable que sea hombre o mujer
		r := rand.IntN(6)
		if slices.Contains([]int{0, 1, 2}, r) {
			m.Gender = "Masculino"
		} else if slices.Contains([]int{3, 4, 5}, r) {
			m.Gender = "Femenino"
		} else {
			m.Gender = "Otro"
		}

		switch m.Gender{
		case "Masculino":
			m.Name = jsonData.MaleFirstNames[rand.IntN(len(jsonData.MaleFirstNames))]
		case "Femenino":
			m.Name = jsonData.FemaleFirstNames[rand.IntN(len(jsonData.FemaleFirstNames))]
		case "Otro":
			random := rand.IntN(2)
			if random == 0 {
				m.Name = jsonData.MaleFirstNames[rand.IntN(len(jsonData.MaleFirstNames))]
			} else {
				m.Name = jsonData.FemaleFirstNames[rand.IntN(len(jsonData.FemaleFirstNames))]
			}		
		}
		m.LastName = jsonData.LastNames[rand.IntN(len(jsonData.LastNames))]
		// año entre 1950 y 2006 (mayor de 18 años)
		year := rand.IntN(56) + 1950
		month := rand.IntN(11) + 1
		var day int
		switch month {
		case 2:
			day = rand.IntN(27) + 1
		case 4, 6, 9, 11:
			day = rand.IntN(29) + 1
		case 1, 3, 5, 7, 8, 10, 12:
			day = rand.IntN(30) + 1
		}
		// creo el dni con la fecha de nacimiento
		m.Dni = createDNI(day, month, year)

		// en la base de datos: '1998-05-22' string
		// consulta: SELECT > CAST('2023-06-26' AS DATE)
		// si lo quiero mostrar en el input lo doy vuelta y listo
		//
		// fijarse bien desp lo del formato fecha
		m.Birthday = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		// m.Birthday = time.Date(year, time.Month(month), day)
		m.MaritalStatus = jsonData.MaritalStatus[rand.IntN(len(jsonData.MaritalStatus))]
		m.Phone = fmt.Sprintf("156%d", rand.IntN(9999999))
		m.Email = fmt.Sprintf("%s%s%s@gmail.com", m.Name, m.LastName, strconv.Itoa(year)[2:])
		m.Address = fmt.Sprintf("%s %d", jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name, rand.IntN(9999))
		m.PostalCode = strconv.Itoa(rand.IntN(8000) + 1000)
		m.District = jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name
		m.MemberNumber = strconv.Itoa(rand.IntN(9999999999))
		cuil := fmt.Sprintf("%d-%s-%d", rand.IntN(9)+20, m.Dni, rand.IntN(8)+1)
		m.Cuil = &cuil
		m.CompanyID = rand.IntN(49) + 1

		switch {
		case year <= 1960:
			m.Category = "Nivel 1: Oficial Múltiple"
		case year <= 1970:
			m.Category = "Nivel 2: Oficial Especializado"
		case year <= 1980:
			m.Category = "Nivel 3: Oficial General"
		case year <= 1990:
			m.Category = "Nivel 4: Medio Oficial"
		case year <= 2000:
			m.Category = "Nivel 5: Ayudante"
		case year <= 2006:
			m.Category = "Nivel 6: Operario Act. Industrial"
		}

		// que sea a los 18 años o mas, entre 18 y 48
		entryYear := rand.IntN(30) + year + 18
		if entryYear > 2024 {
			entryYear = 2024
		}
		entryMonth := rand.IntN(11) + 1
		var entryDay int
		switch month {
		case 2:
			entryDay = rand.IntN(27) + 1
		case 4, 6, 9, 11:
			entryDay = rand.IntN(29) + 1
		case 1, 3, 5, 7, 8, 10, 12:
			entryDay = rand.IntN(30) + 1
		}

		entryDate := time.Date(entryYear, time.Month(entryMonth), entryDay, 0, 0, 0, 0, time.UTC)
		m.EntryDate = &entryDate
		m.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))

		members = append(members, m)
	}

	if len(members) == 0 {
		return nil
	}

	query := "INSERT INTO members (name, last_name, dni, birthday, gender, marital_status, phone, email, address, postal_code, district, member_number, cuil, id_company, category, entry_date, observations) VALUES"
	placeholders := make([]string, 0, len(members))
	args := make([]any, 0, len(members)*17)
	for _, m := range members {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		args = append(args, m.Name, m.LastName, m.Dni, m.Birthday, m.Gender, m.MaritalStatus, m.Phone, m.Email, m.Address, m.PostalCode, m.District, m.MemberNumber, m.Cuil, m.CompanyID, m.Category, m.EntryDate, m.Observations)
	}
	query += strings.Join(placeholders, ",")

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error inserting member: %w", err)
	}
	fmt.Println("members created")
	return nil
}

// func CreateMembers(db *sqlx.DB) error {
// 	type StreetData struct {
// 		Code string
// 		Name string
// 	}
// 	type JSONData struct {
// 		FemaleFirstNames []string
// 		MaleFirstNames   []string
// 		LastNames        []string
// 		MaritalStatus    []string
// 		Genders          []string
// 		Streets          []StreetData
// 	}
//
// 	var m member.Member
//
// 	file, err := os.Open("./internal/bootstrap/creators/data.json")
// 	if err != nil {
// 		fmt.Println("error opening file")
// 		panic(err)
// 	}
// 	decoder := json.NewDecoder(file)
// 	jsonData := JSONData{}
// 	decoder.Decode(&jsonData)
//
// 	for i := 0; i < 100; i++ {
// 		// para que sea mas probable que sea hombre o mujer
// 		r := rand.IntN(6)
// 		if slices.Contains([]int{0, 1, 2}, r) {
// 			m.Gender = "Masculino"
// 		} else if slices.Contains([]int{3, 4, 5}, r) {
// 			m.Gender = "Femenino"
// 		} else {
// 			m.Gender = "Otro"
// 		}
// 		if m.Gender == "Masculino" {
// 			m.Name = jsonData.MaleFirstNames[rand.IntN(len(jsonData.MaleFirstNames))]
// 		} else if m.Gender == "Femenino" {
// 			m.Name = jsonData.FemaleFirstNames[rand.IntN(len(jsonData.FemaleFirstNames))]
// 		} else if m.Gender == "Otro" {
// 			random := rand.IntN(2)
// 			if random == 0 {
// 				m.Name = jsonData.MaleFirstNames[rand.IntN(len(jsonData.MaleFirstNames))]
// 			} else {
// 				m.Name = jsonData.FemaleFirstNames[rand.IntN(len(jsonData.FemaleFirstNames))]
// 			}
// 		}
// 		m.LastName = jsonData.LastNames[rand.IntN(len(jsonData.LastNames))]
// 		// año entre 1950 y 2006 (mayor de 18 años)
// 		year := rand.IntN(56) + 1950
// 		month := rand.IntN(11) + 1
// 		var day int
// 		switch month {
// 		case 2:
// 			day = rand.IntN(27) + 1
// 		case 4, 6, 9, 11:
// 			day = rand.IntN(29) + 1
// 		case 1, 3, 5, 7, 8, 10, 12:
// 			day = rand.IntN(30) + 1
// 		}
// 		// creo el dni con la fecha de nacimiento
// 		m.Dni = createDNI(day, month, year)
//
// 		// en la base de datos: '1998-05-22' string
// 		// consulta: SELECT > CAST('2023-06-26' AS DATE)
// 		// si lo quiero mostrar en el input lo doy vuelta y listo
// 		//
// 		// fijarse bien desp lo del formato fecha
// 		m.Birthday = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
// 		// m.Birthday = time.Date(year, time.Month(month), day)
// 		m.MaritalStatus = jsonData.MaritalStatus[rand.IntN(len(jsonData.MaritalStatus))]
// 		m.Phone = fmt.Sprintf("156%d", rand.IntN(9999999))
// 		m.Email = fmt.Sprintf("%s%s%s@gmail.com", m.Name, m.LastName, strconv.Itoa(year)[2:])
// 		m.Address = fmt.Sprintf("%s %d", jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name, rand.IntN(9999))
// 		m.PostalCode = strconv.Itoa(rand.IntN(8000) + 1000)
// 		m.District = jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name
// 		m.MemberNumber = strconv.Itoa(rand.IntN(9999999999))
// 		cuil := fmt.Sprintf("%d-%s-%d", rand.IntN(9)+20, m.Dni, rand.IntN(8)+1)
// 		m.Cuil = &cuil
// 		m.CompanyID = rand.IntN(49) + 1
//
// 		switch {
// 		case year <= 1960:
// 			m.Category = "Nivel 1: Oficial Múltiple"
// 		case year <= 1970:
// 			m.Category = "Nivel 2: Oficial Especializado"
// 		case year <= 1980:
// 			m.Category = "Nivel 3: Oficial General"
// 		case year <= 1990:
// 			m.Category = "Nivel 4: Medio Oficial"
// 		case year <= 2000:
// 			m.Category = "Nivel 5: Ayudante"
// 		case year <= 2006:
// 			m.Category = "Nivel 6: Operario Act. Industrial"
// 		}
//
// 		// que sea a los 18 años o mas, entre 18 y 48
// 		entryYear := rand.IntN(30) + year + 18
// 		if entryYear > 2024 {
// 			entryYear = 2024
// 		}
// 		entryMonth := rand.IntN(11) + 1
// 		var entryDay int
// 		switch month {
// 		case 2:
// 			entryDay = rand.IntN(27) + 1
// 		case 4, 6, 9, 11:
// 			entryDay = rand.IntN(29) + 1
// 		case 1, 3, 5, 7, 8, 10, 12:
// 			entryDay = rand.IntN(30) + 1
// 		}
//
// 		entryDate := time.Date(entryYear, time.Month(entryMonth), entryDay, 0, 0, 0, 0, time.UTC)
// 		m.EntryDate = &entryDate
// 		m.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))
//
// 		insert, err := db.Query("INSERT INTO members (name, last_name, dni, birthday, gender, marital_status, phone, email, address, postal_code, district, member_number, cuil, id_company, category, entry_date, observations) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", m.Name, m.LastName, m.Dni, m.Birthday, m.Gender, m.MaritalStatus, m.Phone, m.Email, m.Address, m.PostalCode, m.District, m.MemberNumber, m.Cuil, m.CompanyID, m.Category, m.EntryDate, m.Observations)
// 		if err != nil {
// 			return fmt.Errorf("error inserting member: %w", err)
//
// 		}
// 		insert.Close()
// 	}
// 	fmt.Println("members created")
// 	return nil
// }
