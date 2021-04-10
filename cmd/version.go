package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s-secret-injector/pkg/version"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of k8s secret injector",
	Long:  `Print the version of k8s secret injector.`,
	Run: func(cmd *cobra.Command, args []string) {
		k8sSecretInjectorEnvVersion := version.GetVersion()
		fmt.Printf("k8s secret injector version: %v\n", k8sSecretInjectorEnvVersion)
	},
}
