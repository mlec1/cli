package cmd

import (
	"fmt"

	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/spf13/cobra"
)

var nlbDeleteCmd = &cobra.Command{
	Use:     "delete NAME|ID",
	Short:   "Delete a Network Load Balancer",
	Aliases: gRemoveAlias,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmdExitOnUsageError(cmd, "invalid arguments")
		}

		cmdSetZoneFlagFromDefault(cmd)

		return cmdCheckRequiredFlags(cmd, []string{"zone"})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		ctx := exoapi.WithEndpoint(gContext, exoapi.NewReqEndpoint(gCurrentAccount.Environment, zone))
		nlb, err := lookupNLB(ctx, zone, args[0])
		if err != nil {
			return err
		}

		if !force {
			if !askQuestion(fmt.Sprintf("Are you sure you want to delete Network Load Balancer %q?", args[0])) {
				return nil
			}
		}

		if err := cs.DeleteNetworkLoadBalancer(ctx, zone, nlb.ID); err != nil {
			return fmt.Errorf("unable to delete Network Load Balancer: %s", err)
		}

		if !gQuiet {
			cmd.Println("Network Load Balancer deleted successfully")
		}

		return nil
	},
}

func init() {
	nlbDeleteCmd.Flags().BoolP("force", "f", false, cmdFlagForceHelp)
	nlbDeleteCmd.Flags().StringP("zone", "z", "", "Network Load Balancer zone")
	nlbCmd.AddCommand(nlbDeleteCmd)
}
