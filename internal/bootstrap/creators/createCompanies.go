package creators

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/LucasBastino/app-sindicato/internal/features/company"
	"github.com/jmoiron/sqlx"
)

func CreateCompanies(db *sqlx.DB) error {
	type StreetData struct {
		Name string
	}
	type JSONData struct {
		CompaniesNames []string
		Streets         []StreetData
	}

	file, err := os.Open("./internal/bootstrap/creators/data.json")
	if err != nil {
		fmt.Println("error opening json file")
	}
	decoder := json.NewDecoder(file)
	var jsonData JSONData
	decoder.Decode(&jsonData)

	var co company.Company

	// insert, err := database.DB.Query("INSERT INTO companies (Name, CompanyNumber, Address, Contact, CUIT, District, PostalCode, Phone, Observations) VALUES ('SIN EMPRESA', 'POR DEFECTO', 'POR DEFECTO', 'POR DEFECTO', 'POR DEFECTO', 'POR DEFECTO', '0000', 'POR DEFECTO', 'POR DEFECTO')")
	// if err != nil {
	// 	fmt.Println("error inserting enterprise")
	// 	panic(err)
	// }
	// insert.Close()

	for i := 0; i < 50; i++ {
		co.Name = jsonData.CompaniesNames[rand.IntN(len(jsonData.CompaniesNames))]
		companyNumber := strconv.Itoa(rand.IntN(8000) + 1000)
		co.CompanyNumber = &companyNumber
		cuit := fmt.Sprintf("%d-%d-%d", rand.IntN(9)+20, rand.IntN(8999999)+1000000, rand.IntN(8)+1)
		co.Cuit = &cuit
		co.Address = fmt.Sprintf("%s %d", jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name, rand.IntN(9999))
		co.District = jsonData.Streets[rand.IntN(len(jsonData.Streets))].Name
		co.PostalCode = strconv.Itoa(rand.IntN(8000) + 1000)
		co.Phone = fmt.Sprintf("156%d", rand.IntN(9999999))
		co.Contact = fmt.Sprintf("156%d", rand.IntN(9999999))
		co.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))

		insert, err := db.Query("INSERT INTO companies (name, company_number, address, cuit, district, postal_code, phone, contact, observations) VALUES (?,?,?,?,?,?,?,?,?)", co.Name, co.CompanyNumber, co.Address, co.Cuit, co.District, co.PostalCode, co.Phone, co.Contact, co.Observations)
		if err != nil {
			return fmt.Errorf("error inserting companies: %w", err)
		}
		insert.Close()
	}
	fmt.Println("companies created")
	return nil
}
