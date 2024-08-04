package main

import (
	"strconv"

	"github.com/spf13/cobra"
	bladeapiv1alpha1 "github.com/uptime-induestries/compute-blade-agent/api/bladeapi/v1alpha1"
)

func init() {
	cmdFan.AddCommand(cmdFanSetPercent)
	rootCmd.AddCommand(cmdFan)
}

var (
	cmdFan = &cobra.Command{
		Use:   "fan",
		Short: "Fan-related commands for the compute blade",
	}

	cmdFanSetPercent = &cobra.Command{
		Use:     "set-percent <percent>",
		Example: "bladectl fan set-percent 50",
		Short:   "Set the fan speed in percent",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			ctx := cmd.Context()
			client := clientFromContext(ctx)

			// convert string to int
			percent, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			_, err = client.SetFanSpeed(ctx, &bladeapiv1alpha1.SetFanSpeedRequest{
				Percent: int64(percent),
			})

			return err
		},
	}
)
