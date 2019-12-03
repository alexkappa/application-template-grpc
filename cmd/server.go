package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alexkappa/grpc-demo/api"
	"github.com/alexkappa/grpc-demo/api/echo"
	"github.com/alexkappa/grpc-demo/api/health"
	"github.com/alexkappa/grpc-demo/pkg/mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdServer = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve", "run", "start"},
	Short:   "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {

		addr, _ := cmd.Flags().GetString("addr")
		mongoAddr, _ := cmd.Flags().GetString("mongo-addr")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		m, err := mongo.Client(ctx, mongoAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(127)
		}

		server := api.NewServer(
			api.WithAddress(addr),
			api.WithContext(ctx),
			api.WithMongo(m))

		server.Register(
			health.Service(),
			echo.Service())

		return server.Run()
	},
}

func init() {
	cmdRoot.AddCommand(cmdServer)
	cmdServer.Flags().String("addr", ":11001", "address the server listens on")
	cmdServer.Flags().String("mongo-addr", "mongodb://localhost:27017", "mongodb server address")

	viper.BindPFlag("addr", cmdServer.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("addr", cmdServer.PersistentFlags().Lookup("mongo-addr"))
}
