package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xvzf/computeblade-agent/internal/agent"
	"github.com/xvzf/computeblade-agent/pkg/ledengine"
	"github.com/xvzf/computeblade-agent/pkg/log"
	"go.uber.org/zap"
)

func main() {
	var wg sync.WaitGroup

	// setup logger

	zapLogger := zap.Must(zap.NewDevelopment()).With(zap.String("app", "computeblade-agent"))
	_ = zap.ReplaceGlobals(zapLogger.With(zap.String("scope", "global")))
	baseCtx := log.IntoContext(context.Background(), zapLogger)

	ctx, cancelCtx := context.WithCancelCause(baseCtx)
	defer cancelCtx(context.Canceled)

	agent, err := agent.NewComputeBladeAgent(agent.ComputeBladeAgentConfig{
		IdleLedColor:        ledengine.LedColorGreen(0.05),
		IdentifyLedColor:    ledengine.LedColorPurple(0.05),
		CriticalLedColor:    ledengine.LedColorRed(0.3),
		StealthModeEnabled:  false,
		DefaultFanSpeed:     40,
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
		err := agent.Run(ctx)
		if err != nil && err != context.Canceled {
			log.FromContext(ctx).Error("Failed to run agent", zap.Error(err))
			cancelCtx(err)
		}
	}()

	// setup prometheus endpoint
	promHandler := http.NewServeMux()
	promHandler.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: ":9666", Handler: promHandler}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.FromContext(ctx).Error("Failed to start prometheus server", zap.Error(err))
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
			log.FromContext(ctx).Error("Failed to shutdown prometheus server", zap.Error(err))
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
