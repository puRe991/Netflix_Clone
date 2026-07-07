// Command streamflix runs the StreamFlix web server.
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/billing"
	"github.com/pure991/streamflix/internal/config"
	"github.com/pure991/streamflix/internal/db"
	"github.com/pure991/streamflix/internal/httpserver"
	"github.com/pure991/streamflix/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		return err
	}

	st := store.New(pool)
	sessions := auth.NewSessions(cfg.JWTSecret, cfg.AppEnv == "production")
	bill := billing.New(cfg.StripeSecretKey, cfg.StripeBasicPriceID, cfg.StripePremiumPriceID, cfg.StripeWebhookSecret, cfg.AppURL)

	srv, err := httpserver.New(st, sessions, bill, cfg.AppEnv, cfg.AppURL)
	if err != nil {
		return err
	}

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
	}()

	log.Printf("streamflix listening on :%s (env=%s)", cfg.Port, cfg.AppEnv)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
