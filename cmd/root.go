package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRoot = &cobra.Command{
	Use:   "service",
	Short: "A demonstration of gRPC services with gRPC gateway",
}

func Execute() error {
	return cmdRoot.Execute()
}

func init() {
	cmdRoot.PersistentFlags().BoolP("version", "v", false, "print version and exit")
	// cmdRoot.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// cmdRoot.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("version", cmdRoot.PersistentFlags().Lookup("version"))
	// viper.BindPFlag("useViper", cmdRoot.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// cmdRoot.AddCommand(addCmd)
	// cmdRoot.AddCommand(initCmd)
}
