package injector

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	vaultSecretsManager "k8s-secret-injector/pkg/vault"
)

// SanitizedEnviron will hold env without VAULT env vars
type SanitizedEnviron []string

// Vault known Env Vars
var sanitizeEnvmap = map[string]bool{
	"VAULT_TOKEN":                  true,
	"VAULT_ADDR":                   true,
	"VAULT_CACERT":                 true,
	"VAULT_CAPATH":                 true,
	"VAULT_CLIENT_CERT":            true,
	"VAULT_CLIENT_KEY":             true,
	"VAULT_CLIENT_TIMEOUT":         true,
	"VAULT_CLUSTER_ADDR":           true,
	"VAULT_MAX_RETRIES":            true,
	"VAULT_REDIRECT_ADDR":          true,
	"VAULT_SKIP_VERIFY":            true,
	"VAULT_TLS_SERVER_NAME":        true,
	"VAULT_CLI_NO_COLOR":           true,
	"VAULT_RATE_LIMIT":             true,
	"VAULT_NAMESPACE":              true,
	"VAULT_MFA":                    true,
	"VAULT_ROLE":                   true,
	"VAULT_PATH":                   true,
	"VAULT_AUTH_METHOD":            true,
	"VAULT_TRANSIT_KEY_ID":         true,
	"VAULT_TRANSIT_PATH":           true,
	"VAULT_IGNORE_MISSING_SECRETS": true,
	"VAULT_ENV_PASSTHROUGH":        true,
	"VAULT_JSON_LOG":               true,
	"VAULT_LOG_LEVEL":              true,
	"VAULT_REVOKE_TOKEN":           true,
	"VAULT_ENV_DAEMON":             true,
	"VAULT_ENV_FROM_PATH":          true,
}

// Appends variable an entry (name=value) into the environ list.
// VAULT_* variables are not populated into this list.
func (environ *SanitizedEnviron) append(name, value string) {
	if _, ok := sanitizeEnvmap[name]; !ok {
		*environ = append(*environ, fmt.Sprintf("%s=%s", name, value))
	}
}

// InjectSecrets into the sanitized env
func InjectSecrets(secretData map[string]interface{}, environ []string, sanitized SanitizedEnviron) ([]string, error) {
	/*
		go over the current env vars
		if the env var contains a vault: or secret: prefix it will be added to the sanitized env
		if not add all key values from the secret data to the env vars
	*/
	var data map[string]interface{}
	var vaultSecretKey string
	var prefixedEnv bool
	var explicitKey bool

	data = vaultSecretsManager.CastSecretDataToStringMap(secretData)

	for _, env := range environ {
		prefixedEnv = false
		split := strings.SplitN(env, "=", 2)
		name := split[0]
		value := split[1]

		if strings.HasPrefix(value, ">>vault:") {
			value = strings.TrimPrefix(value, ">>")
		}

		hasVaultPrefix := strings.HasPrefix(value, "vault:")
		hasSecretPrefix := strings.HasPrefix(value, "secret:")

		if hasVaultPrefix {
			// API_KEY=secret:API_KEY
			vaultSecretKey = strings.TrimPrefix(value, "vault:")
			prefixedEnv = true
		}

		if hasSecretPrefix {
			vaultSecretKey = strings.TrimPrefix(value, "secret:")
			prefixedEnv = true
		}

		if prefixedEnv {
			// if the secret data contains an explicit key from env add it to the sanitized env
			log.Debugf("Explicit key: %s found in env vars, checking if its in vault secrets...", vaultSecretKey)
			explicitKey = true
			if value, ok := data[vaultSecretKey]; ok {
				log.Debugf("Explicit key: %s found, will be added to the process environment", vaultSecretKey)
				sanitized.append(name, fmt.Sprintf("%v", value))
			} else {
				return nil, fmt.Errorf("Explicit key: %s not found in secrets keys", vaultSecretKey)
			}
		} else {
			// add the env var to the sanitized env
			sanitized.append(name, value)
		}
	}

	if !explicitKey {
		for secretName, secretValue := range data {
			value := fmt.Sprintf("%v", secretValue)
			sanitized.append(secretName, value)
		}
	}
	return sanitized, nil
}
