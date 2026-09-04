package bootstrap

import (
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/functiontemplates"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/http/errorhandler"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)


func InitApp(cfg config.Config, logger logger.Logger) (*fiber.App, func(), error) {
	
	engine := html.New("./internal/views", ".html")
	engine.Reload(true)
	engine.AddFunc("formatAmountAR", functiontemplates.FormatAmountAR)
	engine.AddFunc("formatAmountInput", functiontemplates.FormatAmountInput)
	// engine := html.NewFileSystem(http.FS(viewFiles), ".html")
	// engine := html.NewFileSystem(http.FS(embedfsSub(viewFiles, "src/views")), ".html")

	infra, err := buildInfra(cfg, logger)
	services := buildServices(infra, cfg)
	httpComponents := buildHTTPComponents(services, infra, cfg.Auth.CookieSecure)
	services.license.StartChecker()
	startCron(services.payment, services.backUp, services.idempotency, logger)

	// chequeo si hay tablas creadas, sino creo todo desde cero
	// aca no hace falta usar el adapter, son querys init faciles
	err = initDB(infra.db)
	if err!=nil{
		return nil, nil, fmt.Errorf("failed to initialize database bootstrap: %w", err)
	}

	// todo: esto borrarlo despues
	// err = creators.CreateDemo(infra.db)
	// if err!=nil{
	// 	return nil, nil, fmt.Errorf("failed to create demo registers: %w", err)
	// }

	// Initializing and config app
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 1 * 1024 * 1024, // 1 MiB — enough for HTML forms; no large uploads
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return errorhandler.HandleError(c, err)
		},
	})

	cleanup := func() {
		infra.db.Close()
	}
	
	// Serve static files
	app.Static("/static", "./internal/static")
	// app.StaticFS("/static", http.FS(staticFiles))
	// app.Use("/static", adaptor.HTTPHandler(http.StripPrefix("/static", http.FileServer(http.FS(embedfsSub(staticFiles, "src/static"))))))

	registerMiddlewares(app, httpComponents.auth.middleware, infra.logger, cfg.Auth)
	registerRoutes(app, httpComponents)
	
	return app, cleanup, nil
}
