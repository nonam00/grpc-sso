package app

import (
	grpcapp "grpc-service-ref/internal/app/grpc"
	"grpc-service-ref/internal/services/auth"
	"grpc-service-ref/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
    GRPCSrc *grpcapp.App
}

func New(
    log *slog.Logger,
    grpcPort int,
    connectionString string,
    tokenTTL time.Duration,
) *App {
    //storage, err := sqlite.New(storagePath)
    storage, err := postgres.New(connectionString)
    if err != nil {
        panic(err)
    }

    authService := auth.New(log, storage, storage, tokenTTL)

    grpcApp := grpcapp.New(log, authService, grpcPort)
    
    return &App{
        GRPCSrc: grpcApp,
    }
}
