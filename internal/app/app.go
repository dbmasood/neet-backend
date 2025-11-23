// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evrone/go-clean-template/config"
	amqprpc "github.com/evrone/go-clean-template/internal/controller/amqp_rpc"
	"github.com/evrone/go-clean-template/internal/controller/grpc"
	"github.com/evrone/go-clean-template/internal/controller/http"
	natsrpc "github.com/evrone/go-clean-template/internal/controller/nats_rpc"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/ai"
	"github.com/evrone/go-clean-template/internal/usecase/analytics"
	"github.com/evrone/go-clean-template/internal/usecase/auth"
	"github.com/evrone/go-clean-template/internal/usecase/coupon"
	"github.com/evrone/go-clean-template/internal/usecase/exam"
	"github.com/evrone/go-clean-template/internal/usecase/feed"
	"github.com/evrone/go-clean-template/internal/usecase/leaderboard"
	"github.com/evrone/go-clean-template/internal/usecase/podcast"
	"github.com/evrone/go-clean-template/internal/usecase/practice"
	"github.com/evrone/go-clean-template/internal/usecase/question"
	"github.com/evrone/go-clean-template/internal/usecase/referral"
	"github.com/evrone/go-clean-template/internal/usecase/revision"
	"github.com/evrone/go-clean-template/internal/usecase/translation"
	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/evrone/go-clean-template/internal/usecase/wallet"
	"github.com/evrone/go-clean-template/pkg/grpcserver"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	natsRPCServer "github.com/evrone/go-clean-template/pkg/nats/nats_rpc/server"
	"github.com/evrone/go-clean-template/pkg/postgres"
	rmqRPCServer "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	repos := persistent.New(pg)

	userJWT := jwt.NewService(cfg.JWT.UserSecret, cfg.App.Name, time.Minute*time.Duration(cfg.JWT.TokenTTLMinutes))
	adminJWT := jwt.NewService(cfg.JWT.AdminSecret, cfg.App.Name, time.Minute*time.Duration(cfg.JWT.TokenTTLMinutes))

	// Use-Case
	useCases := usecase.UseCases{
		Auth:        auth.New(repos.User, userJWT),
		User:        user.New(repos.User, repos.Subject, repos.Topic),
		Practice:    practice.New(repos.Practice),
		Revision:    revision.New(repos.Revision),
		Question:    question.New(repos.Question),
		Exam:        exam.New(repos.Exam),
		Podcast:     podcast.New(repos.Podcast),
		Wallet:      wallet.New(repos.Wallet),
		Coupon:      coupon.New(repos.Coupon),
		Referral:    referral.New(repos.Referral),
		AI:          ai.New(repos.AI),
		Analytics:   analytics.New(repos.Analytics),
		Leaderboard: leaderboard.New(repos.Leaderboard),
		Feed:        feed.New(repos.Feed),
	}

	translationUseCase := translation.New(
		repos.Translation,
		webapi.New(),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase, l)

	rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// NATS RPC Server
	natsRouter := natsrpc.NewRouter(translationUseCase, l)

	natsServer, err := natsRPCServer.New(cfg.NATS.URL, cfg.NATS.ServerExchange, natsRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - natsServer - server.New: %w", err))
	}

	// gRPC Server
	grpcServer := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))
	grpc.NewRouter(grpcServer.App, translationUseCase, l)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, translationUseCase, useCases, userJWT, adminJWT, l)

	// Start servers
	rmqServer.Start()
	natsServer.Start()
	grpcServer.Start()
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	case err = <-natsServer.Notify():
		l.Error(fmt.Errorf("app - Run - natsServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

	err = natsServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - natsServer.Shutdown: %w", err))
	}
}
