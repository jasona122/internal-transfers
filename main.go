package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"internal-transfers/internal/api"
	"internal-transfers/internal/config"
	"internal-transfers/internal/repository"
	"internal-transfers/internal/service"

	"internal-transfers/pkg/logger"
)

func main() {
	logger.InitLogger()

	// init configs
	if err := config.LoadEnv(); err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}

	// init db
	dbCfg := config.GetDBConfig()
	db, err := repository.InitDB(dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize db")
	}

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// init services
	accountSvc := service.NewAccountService(accountRepo)
	transactionSvc := service.NewTransactionService(transactionRepo, accountRepo, db)

	// init router
	router := api.NewRouter(accountSvc, transactionSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msg(fmt.Sprintf("Server running on :%s", port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("Server crashed: %v")
	}
}
