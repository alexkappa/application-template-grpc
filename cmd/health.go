package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/alexkappa/service-template-grpc/api/health"
	proto "github.com/alexkappa/service-template-grpc/proto/health/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cmdHealth = &cobra.Command{
	Use:   "health",
	Short: "A health client able to interact with the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		addr, _ := cmd.Flags().GetString("addr")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		con, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer con.Close()

		client := health.Client(con)

		res, err := client.Check(ctx, &proto.HealthCheckRequest{})
		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(res)
	},
}

func init() {
	cmdRoot.AddCommand(cmdHealth)
	cmdHealth.Flags().String("addr", ":11010", "server address")

	viper.BindPFlag("addr", cmdHealth.PersistentFlags().Lookup("addr"))
}
