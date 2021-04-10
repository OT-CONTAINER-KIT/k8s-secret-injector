<div align="left">
    <img src="./static/k8s-secret-injector-logo.svg" height="120" width="120">
</div>

## K8s Secret Injector

k8s-secret-injector is a tool that can connect with multiple secret managers to fetch the secrets and then inject them into the environment variable in a secure way. The motive is to create this injector is to use it with another project of our k8s-vault-webhook but this can be used independently outside Kubernetes as well.

The secret managers which are currently supported:-

- **[Hashicorp Vault](https://www.vaultproject.io/)**

There are some secret managers which are planned to be implemented in future.

- **[AWS Secret Manager](https://aws.amazon.com/secrets-manager/)**
- **[Azure Key Vault](https://azure.microsoft.com/en-in/services/key-vault/)**
- **[GCP Secret Manager](https://cloud.google.com/secret-manager)**

### Supported Features

- k8s-secret-injector can connect with Vault using Kubernetes as backend
- Authenticate with Kubernetes using Serviceaccount mechanism
- Inject secrets directly to the process, i.e. after the injection you cannot read secrets from the environment variable

### Architecture

<div align="center">
    <img src="./static/k8s-secret-injector-arc.png">
</div>

