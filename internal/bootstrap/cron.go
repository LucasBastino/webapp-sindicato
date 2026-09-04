package bootstrap

import (
	"context"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/backup"
	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/robfig/cron/v3"
)

func startCron(paymentService *payment.PaymentService, backupService *backup.BackUpService, idempotencyService *idempotency.IdempotencyService, logger logger.Logger){
	c := cron.New()

	// CREATE PAYMENTS
	// se ejecuta el 1ro de diciembre a las 00:30
    c.AddFunc("30 0 1 12 *", func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := paymentService.CreateYearlyPayments(ctx); err != nil {
			logger.Error("cron error, failed to create yearly payments: ", err)
        }
    })
	// se ejecuta el 10 de diciembre a las 00:30 por si falló el anterior
	 c.AddFunc("30 0 10 12 *", func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := paymentService.CreateYearlyPayments(ctx); err != nil {
			logger.Error("cron error, failed to create yearly payments: ", err)
        }
    })


	// LOGGER
	// se ejecuta todos los dias a las 00:10
	c.AddFunc("10 0 * * *", func() {
		if err := logger.RotateLogFile(); err != nil {
			logger.Error("cron error, failed to rotate log file: ", err)
        }

		if err := logger.CleanupOldLogFiles(); err != nil {
			logger.Error("cron error, failed to ckeanup old log files: ", err)
        }
	})


	// BACKUP
	// se ejecuta todos los dias a las 01:00
	c.AddFunc("0 1 * * *", func() {
		if err := backupService.BackUp(); err != nil {
			logger.Error("cron error, failed to backup: ", err)
		}	
		
		if err := backupService.CleanupOldBackups(); err != nil {
			logger.Error("cron error, failed to cleanup old backups: ", err)
        }
		
	})


	// IDEMPOTENCY KEYS
	// se ejecuta todas las horas
	c.AddFunc("0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

		if err := idempotencyService.CleanupExpiredKeys(ctx); err != nil {
			logger.Error("cron error, failed to cleanup expired idempotency keys: ", err)
        }
	})

    c.Start()
}
