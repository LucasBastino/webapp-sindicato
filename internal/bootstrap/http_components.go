package bootstrap

import (
	"github.com/LucasBastino/webapp-sindicato/internal/auth"
	"github.com/LucasBastino/webapp-sindicato/internal/features/company"
	"github.com/LucasBastino/webapp-sindicato/internal/features/installment"
	"github.com/LucasBastino/webapp-sindicato/internal/features/member"
	"github.com/LucasBastino/webapp-sindicato/internal/features/parent"
	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	"github.com/LucasBastino/webapp-sindicato/internal/features/paymentplan"
	"github.com/LucasBastino/webapp-sindicato/internal/features/user"
	"github.com/LucasBastino/webapp-sindicato/internal/http/pages"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/webapp-sindicato/internal/license"
)


type httpComponents struct {
	auth       	*authModule
	license		*license.LicenseHandler
	idempotency *idempotency.IdempotencyMiddleware
	pages		*pages.PagesHandler

	user		*user.UserHandler

	company 	*company.CompanyHandler
	member     	*member.MemberHandler
	parent     	*parent.ParentHandler

	payment    	*payment.PaymentHandler
	paymentPlan *paymentplan.PaymentPlanHandler
	installment *installment.InstallmentHandler
	

}

type authModule struct {
	handler		*auth.AuthHandler
	middleware	*auth.AuthMiddleware
	service		*auth.AuthService
}

// type idempotencyModule struct {
// 	middleware *idempotency.IdempotencyMiddleware
// 	service    *idempotency.IdempotencyService
// 	repository *idempotency.IdempotencyRepository
// }


func buildHTTPComponents(services *services, app *infra, cookieSecure bool) *httpComponents{
	
	authMiddleware := auth.NewAuthMiddleware(services.auth, services.license)
	authHandler := auth.NewAuthHandler(services.auth, app.logger)
	licenseHandler := license.NewLicenseHandler(services.license)

	userHandler := user.NewUserHandler(services.user, cookieSecure)

	companyHandler := company.NewCompanyHandler(services.company, app.normalizer)
	memberHandler := member.NewMemberHandler(services.member, services.company, app.normalizer)
	parentHandler := parent.NewParentHandler(services.parent, services.member, app.normalizer)

	paymentHandler := payment.NewPaymentHandler(services.payment, app.normalizer)
	paymentPlanHandler := paymentplan.NewPaymentPlanHandler(services.paymentPlan)
	installmentHandler := installment.NewInstallmentHandler(services.installment)
	
	dashboardService := pages.NewDashboardService(services.member, services.company, services.payment, services.paymentPlan)
	pagesHandler := pages.NewPagesHandler(dashboardService)

	idempotencyMiddleware := idempotency.NewIdempotencyMiddleware(services.idempotency)

	authModule := &authModule{
		handler: authHandler,
		service: services.auth,
		middleware: authMiddleware,
	}

	return &httpComponents{
		auth: authModule,
		license: licenseHandler,

		user: userHandler,

		company: companyHandler,
		member: memberHandler,
		parent: parentHandler,

		payment: paymentHandler,
		paymentPlan: paymentPlanHandler,
		installment: installmentHandler,

		pages:  pagesHandler,

		idempotency: idempotencyMiddleware,
	}
}
