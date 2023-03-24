package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ragnar-BY/companies-handler/internal/broker"
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

	pgSettings := postgres.PostgresSettings{
		Addr:     cfg.PostgresAddress,
		Username: cfg.PostgresUser,
		Password: cfg.PostgresPassword,
		Database: cfg.PostgresDB,
	}
	dbClient, err := postgres.NewPostgresClient(pgSettings)
	if err != nil {
		logger.Fatal("can not connect to database", zap.Error(err))
	}
	msgBroker := broker.NewBroker()
	eventSrv := service.NewEventSender(msgBroker)

	companySrv := service.NewCompanyService(dbClient)
	companyUsecase := usecase.NewCompanyUsecase(companySrv, eventSrv)

	userSrv := service.NewUserService(dbClient)
	authSrv := service.NewAuthService([]byte(cfg.JWTKey))
	authUsecase := usecase.NewAuthUsecase(authSrv, userSrv)
	srv := rest.NewServer(cfg.ServerAddress, logger, companyUsecase, authUsecase)

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

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}
	if err := dbClient.Close(); err != nil {
		logger.Fatal("Database service forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")
}
