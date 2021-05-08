package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s-secret-injector/pkg/azure"
)

var (
	azureVaultName string
)

// azureCmd represente the Azure commands
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Fetch secrets from Azure Key Vault",
	Long:  `Fetch secrets from Azure Key Vault`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := &azure.AzureConfig{
			AzureVaultName: azureVaultName,
		}

		secretData := azure.RetrieveSecretFromAzure(*cfg)
		processSecrets(secretData, args)
	},
}

func init() {
	RootCmd.AddCommand(azureCmd)

	viper.SetDefault("azure_vault_name", "test-secret")
	viper.AutomaticEnv()

	azureCmd.Flags().StringVar(&azureVaultName, "azure-vault-name", viper.GetString("azure_vault_name"), "Name of the azure vault (default: test-secret)")
}
