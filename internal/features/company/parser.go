package company

import (
	"strings"

	p "github.com/LucasBastino/app-sindicato/internal/common/utils/parser"
)


func toModel(req request) Company {
	// chequeo los *string por si son nil
	companyNumber := p.StrOrEmpty(req.CompanyNumber)
	cuit := p.StrOrEmpty(req.Cuit)
	// si es un CPA lo paso a mayusculas
	if len(req.PostalCode) == 8{
		req.PostalCode = strings.ToUpper(req.PostalCode)
	}
    return Company{
        Name:            req.Name,
        CompanyNumber:companyNumber,
        Address:         req.Address,
        Cuit:            cuit,
        District:        req.District,
        PostalCode:      req.PostalCode,
        Phone:           req.Phone,
        Contact:         req.Contact,
        Observations:    req.Observations,
        // CreatedAt y UpdatedAt los maneja la DB o el ORM
    }
}

func toResponse(e Company) response {

	// chequeo si los campo del registro obtenido de la db contienen un valor o si son null
	companyNumber := p.StrOrDBNull(e.CompanyNumber)
	cuit := p.StrOrDBNull(e.Cuit)

    return response{
        ID:              e.ID,
        Name:            e.Name,
        CompanyNumber:companyNumber,
        Address:         e.Address,
        Cuit:            cuit,
        District:        e.District,
        PostalCode:      e.PostalCode,
        Phone:           e.Phone,
        Contact:         e.Contact,
        Observations:    e.Observations,
        CreatedAt:       e.CreatedAt.Format("02/01/2006"),
        UpdatedAt:       e.UpdatedAt.Format("02/01/2006"),
		IsDeleted:       e.DeletedAt != nil,
    }
}

// cuando da error el formulario queriendo crear
func toResponseFromRequest(req request) response {
	return response{
		Name:             req.Name,
		CompanyNumber: req.CompanyNumber,
		Address:          req.Address,
		Cuit:             req.Cuit,
		District:         req.District,
		PostalCode:       req.PostalCode,
		Phone:            req.Phone,
		Contact:          req.Contact,
		Observations:     req.Observations,
	}
}

// cuando da error el formulario queriendo editar
func mergetoResponse(e Company, req request) response{

	// chequeo si los campo del registro obtenido de la db contienen un valor o si son null
	companyNumber := p.StrOrDBNull(e.CompanyNumber)
	cuit := p.StrOrDBNull(e.Cuit)
	// de esta manera, puede compararse con un string vacio del request, sino no son del mismo tipo y por lo tanto, no son comparables

	return response{
		Name:             p.MergeField(e.Name, req.Name),
		CompanyNumber: p.MergeField(companyNumber, req.CompanyNumber),
		Address:          p.MergeField(e.Address, req.Address),
		Cuit:             p.MergeField(cuit, req.Cuit),
		District:         p.MergeField(e.District, req.District),
		PostalCode:       p.MergeField(e.PostalCode, req.PostalCode),
		Phone:            p.MergeField(e.Phone, req.Phone),
		Contact:          p.MergeField(e.Contact, req.Contact),
		Observations:     p.MergeField(e.Observations, req.Observations),
		CreatedAt:        e.CreatedAt.Format("02/01/2006"),
		UpdatedAt:        e.UpdatedAt.Format("02/01/2006"),
	}
}

func toOptionResponse(e Company) optionsResponse{
	return optionsResponse{
		ID: e.ID,
		Name: e.Name,
	}
}

func toOptionResponses(companies []Company) []optionsResponse{
	responses := make([]optionsResponse, len(companies))
	for i, company := range companies{
		responses[i] = toOptionResponse(company)
	}
	return responses
}

func toTableResponse(e Company) tableResponse{

	// chequeo si el campo del registro obtenido de la db contienen un valor o si es null
	companyNumber := p.StrOrDBNull(e.CompanyNumber)
	
	return tableResponse{
		ID:            e.ID,
		Name:          e.Name,
		CompanyNumber: companyNumber,
		Address:       e.Address,
		IsDeleted:     e.DeletedAt != nil,
	}
}

func toTableResponses(companies []Company) []tableResponse{
	responses := make([]tableResponse, len(companies))
	for i, company := range companies{
		responses[i] = toTableResponse(company)
	}
	return responses
}

