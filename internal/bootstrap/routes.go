package bootstrap

import (
	"github.com/gofiber/fiber/v2"
)

const (
	RoleEditor = "editor"
	RoleViewer = "viewer"

	ResourceCompany = "company"
	ResourceMember  = "member"
) 

func registerRoutes(app *fiber.App, http *httpComponents) {
	public :=  app.Group("/")

	public.Get("/login", http.auth.handler.RenderLogin)
	public.Post("/login", http.auth.handler.Login)
	public.Get("/expired_session", http.auth.handler.RenderExpiredSession)
	public.Get("/verifyLicense", http.license.CheckUpdatedLicense)
	
	// todo: cambiar despues
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
    return c.SendStatus(fiber.StatusNoContent)
})

	// --------------------------
	
	protected := app.Group("/", http.auth.middleware.VerifyToken)
	protected.Get("/",  http.pages.RenderDashboard)
	// protected.Get("/insufficient_permissions", http.auth.handler.RenderInsufficientPermissions)
	
	
	protected.Post("/logout", http.auth.handler.Logout)
	protected.Get("/users/new",  http.auth.middleware.VerifyAdmin, http.user.RenderAddForm)
	protected.Post("/users",  http.auth.middleware.VerifyAdmin, http.idempotency.VerifyIdempotency, http.auth.handler.Register)
	protected.Get("/user_panel", http.auth.middleware.VerifyAdmin, http.user.RenderPanel)
	protected.Get("/users/:id/permissions", http.auth.middleware.VerifyAdmin, http.user.RenderUpdatePermissionsModal)
	protected.Put("/users/:id/permissions", http.auth.middleware.VerifyAdmin, http.user.UpdatePermissions)
	protected.Get("/users/:id/change-password", http.user.RenderChangePasswordModal)
	protected.Put("/users/:id/change-password", http.user.ChangePassword)
	protected.Delete("/users/:id", http.auth.middleware.VerifyAdmin, http.user.HardDelete)

	protected.Get("/companies/new", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.company.RenderAddForm)
	protected.Get("/companies/:company_id/members/new", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.RenderAddForm)
	protected.Get("/companies/:company_id/members", http.auth.middleware.VerifyRole(ResourceMember, RoleViewer), http.member.RenderTable)
	protected.Get("/companies/:company_id/payments", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.payment.RenderGrid)
	protected.Get("/companies/:company_id/payment_plans/new", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.paymentPlan.RenderAddForm)
	protected.Get("/companies/:company_id/payment_plans", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.paymentPlan.RenderTable)
	protected.Get("/companies/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.company.RenderPage)
	protected.Get("/companies", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.company.RenderTable)
	// protected.Get("/companies/:page", http.company.RenderTable)
	// protected.Get("/companies?view=for_select", http.company.RenderTableForSelect)
	// protected.Delete("/companies/:id/members/:member_id", http.auth.middleware.VerifyDelete, http.company.DeleteMember)
	protected.Put("/companies/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.company.Update)
	protected.Post("/companies", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.idempotency.VerifyIdempotency, http.company.Create)
	protected.Delete("/companies/:id/", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.company.SoftDelete)
	protected.Post("/companies/:id/restore", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.company.Restore)
	protected.Delete("/companies/:id/permanent", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.company.HardDelete)
	// protected.Get("/companies/:id/payments?year=:year", http.payment.RenderTable)
	
	protected.Get("/members/new", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.RenderAddForm)
	protected.Get("/members/:member_id/parents/new", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.parent.RenderAddForm)
	protected.Get("/members/:member_id/parents", http.auth.middleware.VerifyRole(ResourceMember, RoleViewer), http.parent.RenderTable)
	protected.Get("/members/:id", http.auth.middleware.VerifyRole(ResourceMember, RoleViewer), http.member.RenderPage)
	protected.Get("/members", http.auth.middleware.VerifyRole(ResourceMember, RoleViewer), http.member.RenderTable)
	protected.Post("/members", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.idempotency.VerifyIdempotency, http.member.Create)
	protected.Put("/members/:id", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.Update)
	protected.Delete("/members/:id", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.SoftDelete)
	protected.Post("/members/:id/restore", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.Restore)
	protected.Delete("/members/:id/permanent", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.member.HardDelete)
	
	protected.Get("/parents/:id", http.auth.middleware.VerifyRole(ResourceMember, RoleViewer), http.parent.RenderModal)
	protected.Put("/parents/:id", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.parent.Update)
	protected.Post("/members/:member_id/parents", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.idempotency.VerifyIdempotency, http.parent.Create)
	protected.Delete("/parents/:id/permanent", http.auth.middleware.VerifyRole(ResourceMember, RoleEditor), http.parent.HardDelete)
	
	protected.Get("/payments/overdue", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.payment.RenderOverdue)
	protected.Get("/payments/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.payment.RenderModal)
	protected.Put("/payments/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.payment.Update)
	
	protected.Get("/payment_plans/overview", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.paymentPlan.RenderOverview)
	protected.Get("/payment_plans/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.paymentPlan.RenderPage)
	protected.Post("/payment_plans", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.idempotency.VerifyIdempotency, http.paymentPlan.Create)
	protected.Put("/payment_plans/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.paymentPlan.Update)
	protected.Post("/payment_plans/:id/cancel", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.paymentPlan.Cancel)
	protected.Post("/payment_plans/:id/restore", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.paymentPlan.Restore)
	protected.Delete("/payment_plans/:id/permanent", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.paymentPlan.HardDelete)

	protected.Get("/installments/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleViewer), http.installment.RenderModal)
	protected.Put("/installments/:id", http.auth.middleware.VerifyRole(ResourceCompany, RoleEditor), http.installment.Update)

	protected.Get("/dashboard", http.pages.RenderDashboard)
	protected.Get("/reports", http.pages.RenderReports)
	protected.Get("/support", http.pages.RenderSupport)
	protected.Get("/backup_DB", http.backUp.Backup)
}
