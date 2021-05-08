package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

type AzureConfig struct {
	AzureVaultName string
}

// getAzureVaultsClient will return the azure client interface
func getAzureVaultsClient() keyvault.BaseClient {
	vaultsClient := keyvault.New()
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Errorf("Failed to initialize azure auth %v", err)
	}
	vaultsClient.Authorizer = authorizer
	return vaultsClient
}

// getSecret will get the value of secret
func getSecret(secname string, cfg AzureConfig) string {
	basicClient := getAzureVaultsClient()
	secretResp, err := basicClient.GetSecret(context.Background(), "https://"+cfg.AzureVaultName+".vault.azure.net", secname, "")
	if err != nil {
		log.Errorf("unable to get list of secrets: %v", err)
		os.Exit(1)
	}
	return *secretResp.Value
}

// getVault returns an existing vault
func getVault(ctx context.Context, cfg AzureConfig) keyvault.SecretListResultPage {
	vaultsClient := getAzureVaultsClient()
	secretsList, err := vaultsClient.GetSecrets(ctx, "https://"+cfg.AzureVaultName+".vault.azure.net", nil)
	if err != nil {
		log.Errorf("unable to get list of secrets: %v", err)
		os.Exit(1)
	}
	return secretsList
}

// RetrieveSecretFromAzure will retrieve the secret from Azure
func RetrieveSecretFromAzure(cfg AzureConfig) map[string]interface{} {
	secretData := make(map[string]interface{})
	secretList := getVault(context.Background(), cfg)
	for ; secretList.NotDone(); secretList.NextWithContext(context.Background()) {
		secWithoutType := make([]string, 1)
		for _, secret := range secretList.Values() {
			secWithoutType = append(secWithoutType, path.Base(*secret.ID))
		}
		for _, wov := range secWithoutType {
			if wov != "" {
				tempValue := strings.ReplaceAll(wov, "-", "_")
				secretData[strings.ToUpper(tempValue)] = getSecret(wov, cfg)
			}
		}
	}
	return secretData
}
