package main

import (
	"fmt"
	"grpc-service-ref/internal/app"
	"grpc-service-ref/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)

func main() {
    cfg := config.MustLoad()
  
    log := setupLogger(cfg.Env)

    log.Info("starting application", slog.Any("cfg", cfg))
  
    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.PGConn.Host, cfg.PGConn.Port, cfg.PGConn.User, cfg.PGConn.Password, cfg.PGConn.DbName,
    )

    application := app.New(log, cfg.GRPC.Port, psqlInfo, cfg.TokenTTL)

    go application.GRPCSrc.MustRun()
    
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

    sign := <-stop

    log.Info("stopping applicaition", slog.String("signal", sign.String()))

    application.GRPCSrc.Stop()
    
    log.Info("applicaiton stopped")
}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envDev:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envProd:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
      )
    }

    return log;
}
