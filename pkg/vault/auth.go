package vault

import (
	vaultapi "github.com/hashicorp/vault/api"
)

// Client is a Vault client with Kubernetes support
type Client struct {
	Client  *vaultapi.Client
	Logical *vaultapi.Logical
}

// NewClientWithConfig create a new vault client
func NewClientWithConfig(config *vaultapi.Config, vaultCfg *Config) (*Client, error) {
	var clientToken string
	var err error
	rawClient, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	logical := rawClient.Logical()
	client := &Client{Client: rawClient, Logical: logical}

	jwt, err := GetServiceAccountToken(vaultCfg.TokenPath)
	if err != nil {
		return nil, err
	}
	clientToken, err = KubernetesBackendLogin(client, vaultCfg, jwt)
	if err != nil {
		return nil, err
	}

	if err == nil {
		rawClient.SetToken(string(clientToken))
	} else {
		return nil, err
	}
	return client, nil
}
