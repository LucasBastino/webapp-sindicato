package bootstrap

import (
	"github.com/LucasBastino/app-sindicato/internal/auth"
	"github.com/LucasBastino/app-sindicato/internal/backup"
	"github.com/LucasBastino/app-sindicato/internal/features/company"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
	"github.com/LucasBastino/app-sindicato/internal/features/member"
	"github.com/LucasBastino/app-sindicato/internal/features/parent"
	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/LucasBastino/app-sindicato/internal/features/paymentplan"
	"github.com/LucasBastino/app-sindicato/internal/features/user"
	"github.com/LucasBastino/app-sindicato/internal/http/pages"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/app-sindicato/internal/license"
)


type httpComponents struct {
	auth       	*authModule
	license		*license.LicenseHandler
	backUp      *backup.BackUpHandler
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


func buildHTTPComponents(services *services, app *infra) *httpComponents{
	
	authMiddleware := auth.NewAuthMiddleware(services.auth, services.license)
	authHandler := auth.NewAuthHandler(services.auth, services.idempotency, app.logger)
	licenseHandler := license.NewLicenseHandler(services.license)

	userHandler := user.NewUserHandler(services.user)

	companyHandler := company.NewCompanyHandler(services.company, services.idempotency, app.normalizer)
	memberHandler := member.NewMemberHandler(services.member, services.idempotency, app.normalizer)
	parentHandler := parent.NewParentHandler(services.parent, services.idempotency, app.normalizer)

	paymentHandler := payment.NewPaymentHandler(services.payment, app.normalizer)
	paymentPlanHandler := paymentplan.NewPaymentPlanHandler(services.paymentPlan, services.idempotency)
	installmentHandler := installment.NewInstallmentHandler(services.installment)
	
	backUpHandler := backup.NewBackUpHandler(services.backUp, app.logger)
	pagesHandler := pages.NewPagesHandler()

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

		backUp: backUpHandler,
		pages:  pagesHandler,

		idempotency: idempotencyMiddleware,
	}
}
