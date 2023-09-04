package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bladeapiv1alpha1 "github.com/xvzf/computeblade-agent/api/bladeapi/v1alpha1"
	"github.com/xvzf/computeblade-agent/internal/agent"
	"github.com/xvzf/computeblade-agent/pkg/fancontroller"
	"github.com/xvzf/computeblade-agent/pkg/ledengine"
	"github.com/xvzf/computeblade-agent/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup

	// setup logger
	zapLogger := zap.Must(zap.NewDevelopment()).With(zap.String("app", "computeblade-agent"))
	_ = zap.ReplaceGlobals(zapLogger.With(zap.String("scope", "global")))
	baseCtx := log.IntoContext(context.Background(), zapLogger)

	ctx, cancelCtx := context.WithCancelCause(baseCtx)
	defer cancelCtx(context.Canceled)

	computebladeAgent, err := agent.NewComputeBladeAgent(agent.ComputeBladeAgentConfig{
		IdleLedColor:       ledengine.LedColorGreen(0.05),
		IdentifyLedColor:   ledengine.LedColorPurple(0.05),
		CriticalLedColor:   ledengine.LedColorRed(0.3),
		StealthModeEnabled: false,
		FanControllerConfig: fancontroller.FanControllerConfig{
			Steps: []fancontroller.FanControllerStep{
				{Temperature: 40, Speed: 40},
				{Temperature: 55, Speed: 80},
			},
		},
		FanUpdateInterval:   5 * time.Second,
		CriticalTemperature: 60,
	})
	if err != nil {
		log.FromContext(ctx).Error("Failed to create agent", zap.Error(err))
		cancelCtx(err)
	}

	// setup stop signal handlers
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for context cancel or signal
		select {
		case <-ctx.Done():
		case sig := <-sigs:
			// On signal, cancel context
			cancelCtx(fmt.Errorf("signal %s received", sig))
		}
	}()

	// Run agent
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := computebladeAgent.Run(ctx)
		if err != nil && err != context.Canceled {
			log.FromContext(ctx).Error("Failed to run agent", zap.Error(err))
			cancelCtx(err)
		}
	}()

	// Setup GRPC server
	// FIXME add logging middleware
	grpcServer := grpc.NewServer()
	bladeapiv1alpha1.RegisterBladeAgentServiceServer(grpcServer, agent.NewGrpcServiceFor(computebladeAgent))
	wg.Add(1)
	go func() {
		defer wg.Done()
		socketPath := "/tmp/computeblade-agent.sock"
		grpcListen, err := net.Listen("unix", "/tmp/computeblade-agent.sock")
		if err != nil {
			log.FromContext(ctx).Error("Failed to create grpc listener", zap.Error(err))
			cancelCtx(err)
			return
		}
		log.FromContext(ctx).Info("Starting grpc server", zap.String("address", socketPath))
		if err := grpcServer.Serve(grpcListen); err != nil && err != grpc.ErrServerStopped {
			log.FromContext(ctx).Error("Failed to start grpc server", zap.Error(err))
			cancelCtx(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.FromContext(ctx).Info("Shutting down grpc server")
		grpcServer.GracefulStop()
	}()

	// setup prometheus endpoint
	instrumentationHandler := http.NewServeMux()
	instrumentationHandler.Handle("/metrics", promhttp.Handler())
	instrumentationHandler.HandleFunc("/debug/pprof/", pprof.Index)
	instrumentationHandler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	instrumentationHandler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	instrumentationHandler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	instrumentationHandler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := &http.Server{Addr: ":9666", Handler: instrumentationHandler}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.FromContext(ctx).Error("Failed to start prometheus/pprof server", zap.Error(err))
			cancelCtx(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.FromContext(ctx).Error("Failed to shutdown prometheus/pprof server", zap.Error(err))
		}
	}()

	// Wait for context cancel
	wg.Wait()
	if err := ctx.Err(); err != nil && err != context.Canceled {
		log.FromContext(ctx).Fatal("Exiting", zap.Error(err))
	} else {
		log.FromContext(ctx).Info("Exiting")
	}
}
