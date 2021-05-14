package cmd

import (
	"encoding/json"
	"errors"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	vault "k8s-secret-injector/pkg/vault"
)

// secretConfigs in JSON strings format
var (
	secretConfigs             []string
	kubernetesBackend         string
	vaultBackend              string
	tokenPath                 string
	vaultRole                 string
	vaultPath                 string
	secretVersion             string
	vaultUseSecretNamesAsKeys bool
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Fetch and inject secrets from Vault to a given command",
	Long:  `Fetch and inject secrets from Vault to a given command`,
	Args:  validateConfig,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		secretData := make(map[string]interface{}) //nolint
		vaultCfg := &vault.Config{
			Role:              vaultRole,
			TokenPath:         tokenPath,
			Backend:           vaultBackend,
			KubernetesBackend: kubernetesBackend,
		}

		if vaultPath != "" {
			var secretConfig vault.SecretConfigJSON
			secretConfig.Path = vaultPath
			secretConfig.Version = secretVersion
			secretConfig.UseSecretNamesAsKeys = strconv.FormatBool(vaultUseSecretNamesAsKeys)
			secretJSON, _ := json.Marshal(secretConfig)
			secretConfigs = append(secretConfigs, string(secretJSON))
		}

		client, err := vault.NewClientWithConfig(vaultapi.DefaultConfig(), vaultCfg)
		if err != nil {
			exitWithError("Error creating Vault client", err)
		}

		vaultCfg, err = vault.ConfigureVaultSecrets(client.Client, secretConfigs, vaultCfg)
		if err != nil {
			exitWithError("Error configuring Vault paramters", err)
		}

		secretData, err = vault.RetrieveSecrets(client.Client, vaultCfg)
		if err != nil {
			exitWithError("Error retrieving secrets from Vault", err)
		}

		processSecrets(secretData, args)
	},
}

func validateConfig(cmd *cobra.Command, args []string) error {

	if vaultRole == "" {
		return errors.New("Vault role is missing, pass it via --role flag or use VAULT_ROLE environment variable")
	}

	if vaultPath == "" && len(secretConfigs) == 0 {
		return errors.New("Vault secret path is missing  pass it via --path flag, or set VAULT_PATH environment variable,  you can also use --secret-config flag")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(vaultCmd)

	viper.SetDefault("vault_backend", "kubernetes")
	viper.SetDefault("kubernetes_backend", "auth/kubernetes/login")
	viper.SetDefault("vault_role", "")
	viper.SetDefault("token_path", "/var/run/secrets/kubernetes.io/serviceaccount/token")
	viper.SetDefault("vault_path", "")
	viper.SetDefault("vault_secret_version", "")
	viper.SetDefault("vault_use_secret_names_as_keys", false)

	viper.AutomaticEnv()

	// Create flags to variables
	vaultCmd.Flags().StringVarP(&vaultBackend, "backend", "b", viper.GetString("vault_backend"), "Vault authentication backend [kubernetes]")
	vaultCmd.Flags().StringVarP(&kubernetesBackend, "kubernetes-backend", "k", viper.GetString("kubernetes_backend"), "Kubernetes backend authentication path")

	// Role, and Token Path location for kubernetes backend login
	vaultCmd.Flags().StringVar(&vaultRole, "role", viper.GetString("vault_role"), "Vault role (required)")
	vaultCmd.Flags().StringVar(&tokenPath, "token-path", viper.GetString("token_path"), "Kubernetes service account JWT token file path")

	// Single secret
	vaultCmd.Flags().StringVar(&vaultPath, "path", viper.GetString("vault_path"), "Vault secrets path, can be a secret path ending with a \"/\" to get all secrets below that path")
	vaultCmd.Flags().StringVar(&secretVersion, "version", viper.GetString("vault_secret_version"), "Secret version if using a KVv2 (default \"latest\")")
	vaultCmd.Flags().BoolVar(&vaultUseSecretNamesAsKeys, "names-as-keys", viper.GetBool("vault_use_secret_names_as_keys"), "Use secret names as keys (default false)")

	// Multiple secrets via JSON string
	vaultCmd.Flags().StringArrayVarP(
		&secretConfigs,
		"secret-config",
		"",
		[]string{},
		"multiple secrets in JSON string like: '{\"path\": \"/some/secret/path\", \"version\": \"3\", \"use-secret-names-as-keys\":  true}' can be specified a multiple times",
	)
}
