package member

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	p "github.com/LucasBastino/app-sindicato/internal/common/utils/parser"
)


func toModel(req request) (Member, error) {
	companyID, err := strconv.Atoi(req.CompanyID)
	if err!=nil{
		return Member{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse company_id: %w", err), "")
	}

	// ya esta validado, no hace falta chequear el error
	birthdayTime, err := time.Parse("02/01/2006", req.Birthday)
	if err!=nil{
		return Member{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse birthday: %w", err), "")
	}
	entryDateTime, err := time.Parse("02/01/2006", req.EntryDate)
	if err!=nil{
		return Member{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse entry_date: %w", err), "")
	}

	// chequeo el *string por si es nil
	cuil := p.StrOrEmpty(req.Cuil)

	// si es un CPA lo paso a mayusculas
	if len(req.PostalCode) == 8{
		req.PostalCode = strings.ToUpper(req.PostalCode)
	}

	return Member{
        Name:          req.Name,
        LastName:      req.LastName,
        Dni:           req.Dni,
        Birthday:      birthdayTime,
        Gender:        req.Gender,
        MaritalStatus: req.MaritalStatus,
        Phone:         req.Phone,
        Email:         req.Email,
        Address:       req.Address,
        PostalCode:    req.PostalCode,
        District:      req.District,
        MemberNumber:  req.MemberNumber,
        Cuil:          cuil,
        CompanyID:  companyID,
        Category:      req.Category,
        EntryDate:     entryDateTime,
        Observations:  req.Observations,
		CompanyName: req.CompanyName,
        // CreatedAt y UpdatedAt los maneja la DB, no el request
    }, nil

}

func toResponse(m Member) response {

	// chequeo si el campo del registro obtenido de la db contienen un valor o si es null
	cuil := p.StrOrDBNull(m.Cuil)

    return response{
        ID:            m.ID,
        Name:          m.Name,
        LastName:      m.LastName,
        Dni:           m.Dni,
        Birthday:      m.Birthday.Format("02/01/2006"),
        Gender:        m.Gender,
        MaritalStatus: m.MaritalStatus,
        Phone:         m.Phone,
        Email:         m.Email,
        Address:       m.Address,
        PostalCode:    m.PostalCode,
        District:      m.District,
        MemberNumber:  m.MemberNumber,
        // Affiliated:    m.Affiliated,
        Cuil:          cuil,
        CompanyID:  m.CompanyID,
        Category:      m.Category,
        EntryDate:     m.EntryDate.Format("02/01/2006"),
        Observations:  m.Observations,
        CreatedAt:     m.CreatedAt.Format("02/01/2006"),
        UpdatedAt:     m.UpdatedAt.Format("02/01/2006"),
		CompanyName: m.CompanyName,
    }
}

// func toResponses(memberList []Member) []response{
// 	ResponseList := make([]response, len(memberList))
// 	for i, memberModel := range memberList{
// 		ResponseList[i] = toResponse(memberModel)
// 	}
// 	return ResponseList
// }

// cuando da error el formulario queriendo crear
func toResponseFromRequest(req request) (response, error) {
	companyID, err := strconv.Atoi(req.CompanyID)
	if err!=nil{
		return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse company_id: %w", err), "")
	}

	return response{
		Name:          req.Name,
		LastName:      req.LastName,
		Dni:           req.Dni,
		Birthday:      req.Birthday,
		Gender:        req.Gender,
		MaritalStatus: req.MaritalStatus,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		PostalCode:    req.PostalCode,
		District:      req.District,
		MemberNumber:  req.MemberNumber,
		Cuil:          req.Cuil,
		CompanyID:  companyID,
		Category:      req.Category,
		EntryDate:     req.EntryDate,
		Observations:  req.Observations,
		CompanyName: req.CompanyName,
	}, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(m Member, req request) (response, error) {
	companyID, err := strconv.Atoi(req.CompanyID)
	if err!=nil{
		return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse company_id: %w", err), "")
	}
	// chequeo si el campo del registro obtenido de la db contienen un valor o si es null
	cuil := p.StrOrDBNull(m.Cuil)
	// de esta manera, puede compararse con un string vacio del request, sino no son del mismo tipo y por lo tanto, no son comparables

    return response{
    	Name:           p.MergeField(m.Name, req.Name),
		LastName:       p.MergeField(m.LastName, req.LastName),
		Dni:            p.MergeField(m.Dni, req.Dni),
		Birthday:       p.MergeField(m.Birthday.Format("02/01/2006"), req.Birthday),
		Gender:         p.MergeField(m.Gender, req.Gender),
		MaritalStatus:  p.MergeField(m.MaritalStatus, req.MaritalStatus),
		Phone:          p.MergeField(m.Phone, req.Phone),
		Email:          p.MergeField(m.Email, req.Email),
		Address:        p.MergeField(m.Address, req.Address),
		PostalCode:     p.MergeField(m.PostalCode, req.PostalCode),
		District:       p.MergeField(m.District, req.District),
		MemberNumber:   p.MergeField(m.MemberNumber, req.MemberNumber),
		// Affiliated:     affiliated,
		Cuil:           p.MergeField(cuil, req.Cuil),
		CompanyID:   companyID,
		Category:       p.MergeField(m.Category, req.Category),
		EntryDate:      p.MergeField(m.EntryDate.Format("02/01/2006"), req.EntryDate),
		Observations:   p.MergeField(m.Observations, req.Observations),
		CreatedAt:      m.CreatedAt.Format("02/01/2006"),
		UpdatedAt:      m.UpdatedAt.Format("02/01/2006"),
		CompanyName: p.MergeField(m.CompanyName, req.CompanyName),
    }, nil
}

func toTableResponse(m Member) tableResponse {
	isDeleted := m.DeletedAt != nil
	isInactive := !isDeleted && m.CompanyDeletedAt != nil

    return tableResponse{
        ID:          m.ID,
        Name:        m.Name,
        LastName:    m.LastName,
        Dni:         m.Dni,
        CompanyName: m.CompanyName,
        IsDeleted:   isDeleted,
        IsInactive:  isInactive,
    }
}

func toTableResponses(memberList []Member) []tableResponse{
	responses := make([]tableResponse, len(memberList))
	for i, member := range memberList{
		responses[i] = toTableResponse(member)
	}
	return responses
}


/* 
type MemberParser struct{}

func (parser MemberParser) ParseModel(c *fiber.Ctx) (Member, error) {
	m := Member{}
	m.Name = strings.TrimSpace(c.FormValue("name"))
	m.LastName = strings.TrimSpace(c.FormValue("last-name"))
	m.Dni = strings.TrimSpace(c.FormValue("dni"))
	m.Birthday = strings.TrimSpace(c.FormValue("birthday"))
	m.Gender = strings.TrimSpace(c.FormValue("gender"))
	m.MaritalStatus = strings.TrimSpace(c.FormValue("marital-status"))
	m.Phone = strings.TrimSpace(c.FormValue("phone"))
	m.Email = strings.TrimSpace(c.FormValue("email"))
	m.Address = strings.TrimSpace(c.FormValue("address"))
	m.PostalCode = strings.TrimSpace(c.FormValue("postal-code"))
	m.District = strings.TrimSpace(c.FormValue("district"))
	m.MemberNumber = strings.TrimSpace(c.FormValue("member-number"))
	affiliated, err := strconv.ParseBool(strings.TrimSpace(c.FormValue("affiliated")))
	if err != nil {
		customError.InternalError.Msg = err.Error()
		return Member{}, errorHandler.HandleError(c, customError.InternalError)
	}
	m.Affiliated = affiliated
	m.Cuil = strings.TrimSpace(c.FormValue("cuil"))
	CompanyIDStr := strings.TrimSpace(c.FormValue("id-company"))
	if CompanyIDStr == "" {
		m.CompanyID = 0
		// este valor igualmente no se usa
		// es solamente para que no aparezca un error
	} else {
		CompanyID, err := strconv.Atoi(CompanyIDStr)
		if err != nil {
			customError.StrConvError.Msg = err.Error()
			return Member{}, customError.StrConvError
		}
		m.CompanyID = CompanyID
	}
	m.Category = strings.TrimSpace(c.FormValue("category"))
	m.EntryDate = strings.TrimSpace(c.FormValue("entry-date"))
	m.Observations = strings.TrimSpace(c.FormValue("observations"))
	return m, nil
} */