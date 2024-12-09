package cmd

import (
	"fmt"

	"github.com/exoscale/cli/pkg/globalstate"
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/spf13/cobra"
)

func (c *dbaasUserCreateCmd) createMysql(cmd *cobra.Command, _ []string) error {

	ctx := gContext

	client, err := switchClientZoneV3(ctx, globalstate.EgoscaleV3Client, v3.ZoneName(c.Zone))
	if err != nil {
		return err
	}

	req := v3.CreateDBAASMysqlUserRequest{Username: v3.DBAASUserUsername(c.Username)}
	if c.MysqlAuthenticationMethod != "" {
		req.Authentication = v3.EnumMysqlAuthenticationPlugin(c.MysqlAuthenticationMethod)
	}

	op, err := client.CreateDBAASMysqlUser(ctx, c.Name, req)

	if err != nil {
		return err
	}

	decorateAsyncOperation(fmt.Sprintf("Creating DBaaS user %q", c.Username), func() {
		op, err = client.Wait(ctx, op, v3.OperationStateSuccess)
	})

	if err != nil {
		return err
	}

	if !globalstate.Quiet {
		return c.outputFunc((&dbaasUserShowCmd{
			Name:     c.Name,
			Zone:     c.Zone,
			Username: c.Username,
		}).showMysql(ctx))
	}

	return nil
}