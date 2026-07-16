package bootstrap

import (
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/config"
	"github.com/LucasBastino/app-sindicato/internal/http/errorhandler"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)


func InitApp(cfg config.Config, logger logger.Logger) (*fiber.App, func(), error) {
	
	engine := html.New("./internal/views", ".html")
	// engine := html.NewFileSystem(http.FS(viewFiles), ".html")
	// engine := html.NewFileSystem(http.FS(embedfsSub(viewFiles, "src/views")), ".html")

	infra, err := buildInfra(cfg, logger)
	services := buildServices(infra, cfg)
	httpComponents := buildHTTPComponents(services, infra)
	services.license.StartChecker()
	startCron(services.payment, services.installment, services.backUp, services.idempotency, logger)

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
		Views: engine,
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

	registerMiddlewares(app, httpComponents.auth.middleware, infra.logger)
	registerRoutes(app, httpComponents)
	
	return app, cleanup, nil
}
