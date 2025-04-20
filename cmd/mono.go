package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/internal/config"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	appjwt "github.com/namhq1989/ezfeedback-api-server/internal/jwt"
	"github.com/namhq1989/ezfeedback-api-server/internal/mailer"
	"github.com/namhq1989/ezfeedback-api-server/internal/monitoring"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/waiter"
	"github.com/namhq1989/go-utilities/appcontext"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type app struct {
	cfg        config.Server
	database   *database.Database
	caching    *caching.Caching
	jwt        *appjwt.JWT
	queue      *queue.Queue
	mailer     *mailer.Mailer
	monitoring *monitoring.Monitoring
	rest       *echo.Echo
	rpc        *grpc.Server
	waiter     waiter.Waiter
	modules    []monolith.Module
}

func (a *app) Config() config.Server {
	return a.cfg
}

func (a *app) Database() *database.Database {
	return a.database
}

func (a *app) Rest() *echo.Echo {
	return a.rest
}

func (a *app) RPC() *grpc.Server {
	return a.rpc
}

func (a *app) Waiter() waiter.Waiter {
	return a.waiter
}

func (a *app) Caching() *caching.Caching {
	return a.caching
}

func (a *app) Monitoring() *monitoring.Monitoring {
	return a.monitoring
}

func (a *app) JWT() *appjwt.JWT {
	return a.jwt
}

func (a *app) Queue() *queue.Queue {
	return a.queue
}

func (a *app) Mailer() *mailer.Mailer {
	return a.mailer
}

func (a *app) startupModules() error {
	ctx := appcontext.NewRest(a.Waiter().Context())

	for _, module := range a.modules {
		if err := module.Startup(ctx, a); err != nil {
			return err
		} else {
			fmt.Printf("🚀 module %s started\n", module.Name())
		}
	}

	return nil
}

func (a *app) waitForRest(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("*** api server started", a.cfg.RestPort)
		defer fmt.Println("*** api server shutdown")

		if err := a.rest.Start(a.cfg.RestPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.rest.Logger.Fatal("shutting down the server")
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("*** api server to be shutdown")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.rest.Shutdown(timeoutCtx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (a *app) waitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.GRPCPort)
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("*** rpc server started", a.cfg.GRPCPort)
		defer fmt.Println("*** rpc server shutdown")
		if err = a.RPC().Serve(listener); err != nil && !errors.Is(grpc.ErrServerStopped, err) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("*** rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			a.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(30 * time.Second)
		select {
		case <-timeout.C:
			// Force it to stop
			a.RPC().Stop()
			return fmt.Errorf("*** rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}
