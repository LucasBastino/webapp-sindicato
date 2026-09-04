package pages

import (
	"context"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/features/company"
	"github.com/LucasBastino/webapp-sindicato/internal/features/member"
	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	"github.com/LucasBastino/webapp-sindicato/internal/features/paymentplan"
)

type DashboardService struct {
	members     *member.MemberService
	companies   *company.CompanyService
	payments    *payment.PaymentService
	paymentPlan *paymentplan.PaymentPlanService
}

func NewDashboardService(
	members *member.MemberService,
	companies *company.CompanyService,
	payments *payment.PaymentService,
	paymentPlan *paymentplan.PaymentPlanService,
) *DashboardService {
	return &DashboardService{
		members:     members,
		companies:   companies,
		payments:    payments,
		paymentPlan: paymentPlan,
	}
}

type DashboardSummary struct {
	ActiveMembers         int
	ActiveCompanies       int
	OverduePaymentsCount  int
	OverduePaymentsAmount float32
	PaymentPlansCount     int
	RecentMembers         []RecentMemberItem
	RecentCompanies       []RecentCompanyItem
}

type RecentMemberItem struct {
	ID          int
	DisplayName string
	CompanyName string
	CreatedAt   string
}

type RecentCompanyItem struct {
	ID          int
	Name        string
	MemberCount int
}

func (s *DashboardService) GetSummary(ctx context.Context) (DashboardSummary, error) {
	activeMembers, err := s.members.CountActive(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}

	activeCompanies, err := s.companies.CountActive(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}

	overdue, err := s.payments.CountOverdueSummary(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}

	plansCount, err := s.paymentPlan.CountAll(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}

	recentMembers, err := s.members.FindRecent(ctx, 5)
	if err != nil {
		return DashboardSummary{}, err
	}

	recentCompanies, err := s.companies.FindRecent(ctx, 5)
	if err != nil {
		return DashboardSummary{}, err
	}

	memberItems := make([]RecentMemberItem, 0, len(recentMembers))
	for _, m := range recentMembers {
		memberItems = append(memberItems, RecentMemberItem{
			ID:          m.ID,
			DisplayName: fmt.Sprintf("%s, %s", m.LastName, m.Name),
			CompanyName: m.CompanyName,
			CreatedAt:   m.CreatedAt.Format("02/01/2006"),
		})
	}

	companyItems := make([]RecentCompanyItem, 0, len(recentCompanies))
	for _, c := range recentCompanies {
		companyItems = append(companyItems, RecentCompanyItem{
			ID:          c.ID,
			Name:        c.Name,
			MemberCount: c.MemberCount,
		})
	}

	return DashboardSummary{
		ActiveMembers:         activeMembers,
		ActiveCompanies:       activeCompanies,
		OverduePaymentsCount:  overdue.Count,
		OverduePaymentsAmount: overdue.Amount,
		PaymentPlansCount:     plansCount,
		RecentMembers:         memberItems,
		RecentCompanies:       companyItems,
	}, nil
}
