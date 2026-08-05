package bootstrap

import (
	"time"

	"github.com/LucasBastino/app-sindicato/internal/auth"
	"github.com/LucasBastino/app-sindicato/internal/backup"
	"github.com/LucasBastino/app-sindicato/internal/config"
	"github.com/LucasBastino/app-sindicato/internal/features/company"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
	"github.com/LucasBastino/app-sindicato/internal/features/member"
	"github.com/LucasBastino/app-sindicato/internal/features/parent"
	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/LucasBastino/app-sindicato/internal/features/paymentplan"
	"github.com/LucasBastino/app-sindicato/internal/features/user"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/app-sindicato/internal/license"
)

type services struct {
	auth    	*auth.AuthService
	license 	*license.LicenseService
	idempotency *idempotency.IdempotencyService

	user 		*user.UserService

	company		*company.CompanyService
	member  	*member.MemberService
	parent  	*parent.ParentService

	payment     *payment.PaymentService
	paymentPlan *paymentplan.PaymentPlanService
	installment *installment.InstallmentService


	backUp 		*backup.BackUpService
}

func buildServices(infra *infra, cfg config.Config) *services {
	licenseService := license.NewLicenseService(infra.logger, cfg.License)

	idempotencyRepo := idempotency.NewIdempotencyRepository(infra.db)
	idempotencyService := idempotency.NewIdempotencyService(idempotencyRepo, 24*time.Hour)

	userRepo := user.NewUserRepository(infra.db)
	userService := user.NewUserService(userRepo, infra.hasher, idempotencyService)

	authRepo := auth.NewAuthRepository(infra.db)
	authService := auth.NewAuthService(authRepo, userService, infra.hasher, infra.tokenGen, cfg.Auth)

	parentRepo := parent.NewParentRepository(infra.db)
	parentService := parent.NewParentService(parentRepo, idempotencyService)

	memberRepo := member.NewMemberRepository(infra.db)
	memberService := member.NewMemberService(memberRepo, idempotencyService)

	paymentRepo := payment.NewPaymentRepository(infra.db)
	paymentService := payment.NewPaymentService(paymentRepo)

	installmentRepo := installment.NewInstallmentRepository(infra.db)
	installmentService := installment.NewInstallmentService(installmentRepo)

	paymentPlanRepo := paymentplan.NewPaymentPlanRepository(infra.db)
	paymentPlanService := paymentplan.NewPaymentPlanService(paymentPlanRepo, paymentRepo, installmentRepo, idempotencyService, infra.logger)
	installmentService.SetPaymentPlanStatusRefresher(paymentPlanService)

	companyRepo := company.NewCompanyRepository(infra.db)
	companyService := company.NewCompanyService(companyRepo, paymentService, idempotencyService)
	paymentService.SetCompanyReader(companyService)

	backUpRepo := backup.NewBackUpRepository(infra.db)
	backUpService := backup.NewBackUpService(backUpRepo)

	return &services{
		auth:	    	authService,
		license: 		licenseService,

		user: 			userService,

		company: 		companyService,
		member:  		memberService,
		parent: 		parentService,

		payment:    	paymentService,
		paymentPlan: 	paymentPlanService,
		installment:	installmentService,

		backUp: 		backUpService,
		idempotency:	idempotencyService,
	}
}
