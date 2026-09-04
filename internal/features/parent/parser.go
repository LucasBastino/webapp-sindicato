package parent

import (
	"fmt"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/webapp-sindicato/internal/common/utils/parser"
	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)


func toModel(req request) (Parent, error) {
	var birthdayTime time.Time
	if req.Birthday != "" {
		parsed, err := v.ParseDMY(req.Birthday)
		if err != nil {
			return Parent{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse birthday: %w", err), "")
		}
		birthdayTime = parsed
	}

	cuil := pu.StrOrEmpty(req.Cuil)

	return Parent{
		Name:         req.Name,
		LastName:     req.LastName,
		Relationship: req.Relationship,
		Birthday:     birthdayTime,
		Gender:       req.Gender,
		Cuil:         cuil,
		Observations: req.Observations,
	}, nil
}

// cuando da error el formulario queriendo crear
func toResponseFromRequest(req request) (response, error) {

	return response{
		Name:         req.Name,
		LastName:     req.LastName,
		Relationship: req.Relationship,
		Birthday:     req.Birthday,
		Gender:       req.Gender,
		Cuil:         req.Cuil,
		Observations: req.Observations,
	}, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(p Parent, req request) (response, error) {

	// chequeo si el campo del registro obtenido de la db contienen un valor o si es null
	cuil := pu.StrOrDBNull(p.Cuil)
	dbBirthday := ""
	if !p.Birthday.IsZero() {
		dbBirthday = p.Birthday.Format("02/01/2006")
	}

	return response{
		ID:           p.ID,
		Name:         pu.MergeField(p.Name, req.Name),
		LastName:     pu.MergeField(p.LastName, req.LastName),
		Relationship: pu.MergeField(p.Relationship, req.Relationship),
		Birthday:     pu.MergeField(dbBirthday, req.Birthday),
		Gender:       pu.MergeField(p.Gender, req.Gender),
		Cuil:         pu.MergeField(cuil, req.Cuil),
		Observations: pu.MergeField(p.Observations, req.Observations),
	}, nil
}

func toResponse(p Parent) response {

	// chequeo si el campo del registro obtenido de la db contienen un valor o si es null
	cuil := pu.StrOrDBNull(p.Cuil)
	birthday := ""
	if !p.Birthday.IsZero() {
		birthday = p.Birthday.Format("02/01/2006")
	}

	return response{
		ID:           p.ID,
		Name:         p.Name,
		LastName:     p.LastName,
		Relationship: p.Relationship,
		Birthday:     birthday,
		Gender:       p.Gender,
		Cuil:         cuil,
		Observations: p.Observations,
		CreatedAt:    p.CreatedAt.Format("02/01/2006"),
		UpdatedAt:    p.UpdatedAt.Format("02/01/2006"),
	}
}

func toTableResponses(parents []Parent) []response{
	responses := make([]response, len(parents))
	for i, parentModel := range parents{
		responses[i] = toResponse(parentModel)
	}
	return responses
}





/* 
type ParentParser struct{}

func (parser ParentParser) ParseModel(c *fiber.Ctx) (Parent, error) {
	p := Parent{}
	p.Name = strings.TrimSpace(c.FormValue("name"))
	p.LastName = strings.TrimSpace(c.FormValue("last-name"))
	p.Rel = strings.TrimSpace(c.FormValue("rel"))
	p.Gender = strings.TrimSpace(c.FormValue("gender"))
	p.Birthday = strings.TrimSpace(c.FormValue("birthday"))
	p.Cuil = strings.TrimSpace(c.FormValue("cuil"))
	MemberIDStr := strings.TrimSpace(c.FormValue("id-member"))
	MemberID, err := strconv.Atoi(MemberIDStr)
	if err != nil {
		customError.StrConvError.Msg = err.Error()
		return Parent{}, customError.StrConvError
	}
	p.MemberID = MemberID

	return p, nil
}
 */