package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/alexkappa/service-template-grpc/api/echo"
	proto "github.com/alexkappa/service-template-grpc/proto/echo/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cmdEcho = &cobra.Command{
	Use:   "echo",
	Short: "An echo client able to interact with the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		addr, _ := cmd.Flags().GetString("addr")
		value, _ := cmd.Flags().GetString("value")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		con, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer con.Close()

		client := echo.Client(con)

		res, err := client.Echo(ctx, &proto.EchoRequest{Value: value})
		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(res)
	},
}

func init() {
	cmdRoot.AddCommand(cmdEcho)
	cmdEcho.Flags().String("addr", ":11010", "server address")
	cmdEcho.Flags().String("value", "", "value to send")

	viper.BindPFlag("addr", cmdEcho.PersistentFlags().Lookup("addr"))
}
