package app

import (
	  "log/slog"
	  grpcapp "grpc-service-ref/internal/app/grpc"
	  "time"
)

type App struct {
    GRPCSrc *grpcapp.App
}

func New(
    log *slog.Logger,
    grpcPort int,
    storagePath string,
    tokenTTL time.Duration,
) *App {
    grpcApp := grpcapp.New(log, grpcPort)
    
    return &App{
        GRPCSrc: grpcApp,
    }
}
