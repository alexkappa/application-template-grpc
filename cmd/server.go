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
	Short:   "Starts the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		httpAddr, _ := cmd.Flags().GetString("http-addr")
		rpcAddr, _ := cmd.Flags().GetString("rpc-addr")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		server := api.NewServer(
			api.WithHTTPAddress(httpAddr),
			api.WithRPCAddress(rpcAddr),
			api.WithContext(ctx))

		server.Register(
			health.Service(),
			echo.Service(store.NewInMemoryKVStore()))

		return server.Serve()
	},
}

func init() {
	cmdRoot.AddCommand(cmdServer)
	cmdServer.Flags().String("http-addr", ":11001", "address the http server listens on")
	cmdServer.Flags().String("rpc-addr", ":11010", "address the rpc server listens on")

	viper.BindPFlag("addr", cmdServer.PersistentFlags().Lookup("addr"))
}
