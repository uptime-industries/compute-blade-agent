package main

import (
	"github.com/spf13/cobra"
	bladeapiv1alpha1 "github.com/xvzf/computeblade-agent/api/bladeapi/v1alpha1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	cmdIdentify.Flags().Bool("confirm", false, "confirm the identify state")
	cmdIdentify.Flags().Bool("wait", false, "Wait for the identify state to be confirmed (e.g. by a physical button press)")
	rootCmd.AddCommand(cmdIdentify)
}

var cmdIdentify = &cobra.Command{
	Use:   "identify",
	Short: "interact with the compute-blade identity LED",
	RunE: runIdentity,
}

func runIdentity(cmd *cobra.Command, _ []string) error {
	var err error

	ctx := cmd.Context()
	client := clientFromContext(ctx)

	// Get flags
	confirm, err := cmd.Flags().GetBool("confirm")
	if err != nil {
		return err
	}
	wait, err := cmd.Flags().GetBool("wait")
	if err != nil {
		return err
	}

	// Check if we should wait for the identify state to be confirmed
	event := bladeapiv1alpha1.Event_IDENTIFY
	if confirm {
		event = bladeapiv1alpha1.Event_IDENTIFY_CONFIRM
	}

	// Emit the event to the computeblade-agent
	_, err = client.EmitEvent(ctx, &bladeapiv1alpha1.EmitEventRequest{Event: event})
	if err != nil {
		return err
	}

	// Check if we should wait for the identify state to be confirmed
	if wait {
		_, err := client.WaitForIdentifyConfirm(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}
	}

	return nil
}
