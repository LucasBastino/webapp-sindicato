package bootstrap

import (
	"fmt"

	authadapters "github.com/LucasBastino/webapp-sindicato/internal/auth/adapters"
	authports "github.com/LucasBastino/webapp-sindicato/internal/auth/ports"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/database"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/LucasBastino/webapp-sindicato/internal/security/password"
	"github.com/jmoiron/sqlx"
)

type infra struct {
	db       	  	*sqlx.DB
	hasher    	 	password.Hasher
	logger    	 	logger.Logger
	tokenGen  		authports.TokenGenerator
	normalizer		authports.ClaimsNormalizer
}

func buildInfra(cfg config.Config, logger logger.Logger) (*infra, error){
	
	// Database connection
	db, err := database.Connect(cfg.Database)
	if err!=nil{
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}


	// Initializing hasher, logger, normalizer and token generator
	hasher := password.NewBcryptHasher()
	normalizer := authadapters.NewJwtNormalizer()
	tokenGen := authadapters.NewJwtTokenGenerator(normalizer)

	// al hacer var log logger.Logger obligo a que log solo pueda usar los metodos definidos por la interfaz Logger
	// no es lo mismo si hago:
	// log := newSlogger()
	// ahi estoy usando la implementacion concreta, un slog logger, no la interfaz
	// por lo que no obligo a "log" a usar solo los metodos definidos por Logger
	// por lo tanto, puedo usar otros especificos del struct y ya estoy rompiendo con el desacoplamiento
	// si quiero cambiar de logger lo unico que tengo que hacer es logger.NewZapLogger() por ejemplo
	
	return &infra{
		db: db,
		hasher: hasher,
		logger: logger,
		tokenGen: tokenGen,
		normalizer: normalizer,
	}, nil
}