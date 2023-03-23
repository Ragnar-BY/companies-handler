package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ragnar-BY/companies-handler/internal/config"
	"github.com/Ragnar-BY/companies-handler/internal/controllers/rest"
	"github.com/Ragnar-BY/companies-handler/internal/repository/postgres"
	"github.com/Ragnar-BY/companies-handler/internal/service"
	"github.com/Ragnar-BY/companies-handler/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatal("can not load config", zap.Error(err))
	}

	dbClient, err := postgres.NewPostgresClient(postgres.PostgresSettings{
		Addr:     cfg.PostgreSQLAddr,
		Username: cfg.PostgreSQLUser,
		Password: cfg.PostgreSQLPassword,
		Database: cfg.PostgreSQLDatabase,
	})
	if err != nil {
		logger.Fatal("can not connect to database", zap.Error(err))
	}
	companySrv := service.NewCompanyService(dbClient)
	companyUsecase := usecase.NewCompanyUsecase(companySrv)

	userSrv := service.NewUserService(dbClient)
	userUsecase := usecase.NewUserUsecase(userSrv)
	srv := rest.NewServer(cfg.ServerAddress, logger, companyUsecase, userUsecase)

	go func() {
		err = srv.Run()
		if err != nil {
			logger.Error("can not run server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
