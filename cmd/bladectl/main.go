package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	bladeapiv1alpha1 "github.com/xvzf/computeblade-agent/api/bladeapi/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClientContextKey int

const (
	defaultGrpcClientContextKey     grpcClientContextKey = 0
	defaultGrpcClientConnContextKey grpcClientContextKey = 1
)

var (
	grpcAddr string
	timeout  time.Duration
)

func init() {
	rootCmd.PersistentFlags().
		StringVar(&grpcAddr, "addr", "unix:///tmp/computeblade-agent.sock", "address of the computeblade-agent gRPC server")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", time.Minute, "timeout for gRPC requests")
}

func clientIntoContext(ctx context.Context, client bladeapiv1alpha1.BladeAgentServiceClient) context.Context {
	return context.WithValue(ctx, defaultGrpcClientContextKey, client)
}

func clientFromContext(ctx context.Context) bladeapiv1alpha1.BladeAgentServiceClient {
	client, ok := ctx.Value(defaultGrpcClientContextKey).(bladeapiv1alpha1.BladeAgentServiceClient)
	if !ok {
		panic("grpc client not found in context")
	}
	return client
}

var rootCmd = &cobra.Command{
	Use:   "bladectl",
	Short: "bladectl interacts with the computeblade-agent and allows you to manage hardware-features of your compute blade(s)",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		origCtx := cmd.Context()

		// setup signal handlers for SIGINT and SIGTERM
		ctx, cancelCtx := context.WithTimeout(origCtx, timeout)

		// setup signal handler channels
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			// Wait for context cancel or signal
			select {
			case <-ctx.Done():
			case <-sigs:
				// On signal, cancel context
				cancelCtx()
			}
		}()

		conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to dial grpc server: %w", err)
		}
		client := bladeapiv1alpha1.NewBladeAgentServiceClient(conn)

		cmd.SetContext( clientIntoContext(ctx, client))
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
