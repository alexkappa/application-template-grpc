package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alexkappa/service-template-grpc/api"
	"github.com/alexkappa/service-template-grpc/api/echo"
	"github.com/alexkappa/service-template-grpc/api/health"
	"github.com/alexkappa/service-template-grpc/pkg/store"
)

var cmdServer = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve", "run", "start"},
	Short:   "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {

		addr, _ := cmd.Flags().GetString("addr")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		server := api.NewServer(
			api.WithAddress(addr),
			api.WithContext(ctx))

		server.Register(
			health.Service(),
			echo.Service(store.NewInMemoryEchoStore()))

		return server.Run()
	},
}

func init() {
	cmdRoot.AddCommand(cmdServer)
	cmdServer.Flags().String("addr", ":11001", "address the server listens on")

	viper.BindPFlag("addr", cmdServer.PersistentFlags().Lookup("addr"))
}
