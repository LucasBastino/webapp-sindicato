package payment

import (
	"context"
	"errors"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/jmoiron/sqlx"
)

type PaymentService struct {
	repo *PaymentRepository

	companyReader companyReader
}

func NewPaymentService(repo *PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) SetCompanyReader(reader companyReader) {
	s.companyReader = reader
}

func (s *PaymentService) GetCompanyName(ctx context.Context, id int) (string, error) {
	if s.companyReader == nil {
		return "", apperrors.NewInternalError(errors.New("company reader not configured"), "")
	}
	name, err := s.companyReader.GetName(ctx, id)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *PaymentService) Get(ctx context.Context, id int) (*Payment, error) {
	payment, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if payment == nil{
		return nil, apperrors.NewNotFoundError(errors.New("payment not found"), "")
	}

	return payment, nil
}

func (s *PaymentService) List(ctx context.Context, companyID int, year int) ([]Payment, error) {
	payments, err := s.repo.FindAll(ctx, companyID, year)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return payments, nil
}

func (s *PaymentService) ListPaymentYears(ctx context.Context, companyID int) ([]int, error) {
	years, err := s.repo.GetPaymentYears(ctx, companyID)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return years, nil
}

func (s *PaymentService) Count(ctx context.Context, companyID int) (int, error) {
	count, err := s.repo.Count(ctx, companyID)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return count, nil
}

func (s *PaymentService) CountOverdueSummary(ctx context.Context) (OverdueSummary, error) {
	summary, err := s.repo.CountOverdueSummary(ctx)
	if err != nil {
		return OverdueSummary{}, apperrors.NewDatabaseError(err, "")
	}
	return summary, nil
}

func (s *PaymentService) ListOverdue(ctx context.Context) ([]OverdueCompanyGroup, error) {
	rows, err := s.repo.FindAllOverdue(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	if len(rows) == 0 {
		return nil, nil
	}

	byCompany := make(map[int]*OverdueCompanyGroup)
	order := make([]int, 0)
	for _, row := range rows {
		group, ok := byCompany[row.CompanyID]
		if !ok {
			group = &OverdueCompanyGroup{
				CompanyID:   row.CompanyID,
				CompanyName: row.CompanyName,
				Payments:    make([]overduePaymentResponse, 0),
			}
			byCompany[row.CompanyID] = group
			order = append(order, row.CompanyID)
		}
		amount := amountValue(row.Amount)
		group.TotalAmount += amount
		group.Payments = append(group.Payments, overduePaymentResponse{
			ID:        row.ID,
			MonthName: monthName(row.Month),
			Year:      row.Year,
			Amount:    amount,
			DueDate:   row.DueDate.Format("02/01/2006"),
		})
	}

	groups := make([]OverdueCompanyGroup, 0, len(order))
	for _, id := range order {
		groups = append(groups, *byCompany[id])
	}
	return groups, nil
}


// crea los pagos de todos los meses del año
func (s *PaymentService) CreateYearlyPayments(ctx context.Context) error {
    ids, err := s.companyReader.ListActiveIDs(ctx)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	// creo los pagos de todos los meses del año siguiente
	fromMonth := 1
	year := time.Now().Year() + 1
	// prealoco el slice para tener mayor rendimiento
	payments := make([]Payment, 0, len(ids)*(12-fromMonth+1))
	for _, id := range ids{
		payments = append(payments, buildPayments(id, fromMonth, year)...)
	}
    err = s.repo.BulkInsert(ctx, nil, payments)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

// crea los pagos de los meses restantes
func (s *PaymentService) CreateRemainingPaymentsByID(ctx context.Context, tx *sqlx.Tx, companyID int) error {
	// creo los pagos restantes del año actual desde el mes actual
	fromMonth := int(time.Now().Month())
	year := int(time.Now().Year())
	// prealoco el slice para tener mayor rendimiento
	payments := buildPayments(companyID, fromMonth, year)
	// si estoy en diciembre, creo los pagos del año siguiente tambien,
	// porque el cron ya se ejecuto y la empresa no estaba creada
	if fromMonth == 12{
		payments = append(payments, buildPayments(companyID, 1, year+1)...)
	}
    err := s.repo.BulkInsert(ctx, tx, payments)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}


func (s *PaymentService) Update(ctx context.Context, id int, payment Payment) error {
	paymentDB, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if paymentDB == nil{
		return apperrors.NewNotFoundError(errors.New("payment not found"), "")
	}
	
	if paymentDB.IsInPaymentPlan{
		return apperrors.NewBusinessError(errors.New("failed to update payment: can't update payment while it is in payment plan"), "No puedes editar un pago que está dentro de un plan de pago")
	}

	err = s.repo.Update(ctx, id, payment)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}


func buildPayments(companyID, fromMonth, year int) ([]Payment){
	payments := make([]Payment, 0, 12)
	for month := fromMonth; month <= 12; month++ {
		payments = append(payments, Payment{
			CompanyID: companyID,
			Month:     month,
			Year:      year,
			DueDate:   time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC),
		})
	}
	if len(payments) == 0{
		return nil
	}
	return payments
}


// func (s *PaymentService) ValidateInDB(ctx context.Context, companyID int) map[string]string {
// 	return s.repo.ValidateInDB(ctx, companyID)
// }

// func (s *PaymentService) Create(ctx context.Context, model Payment) (int64, error) {
// 	return s.repo.Insert(ctx, model)
// }