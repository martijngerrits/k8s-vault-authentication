# Kubernetes Vault authentication

### Basic example
```go
import (
    k8sva "github.com/martijngerrits/k8s-vault-authentication"
    "log"
)

client, err := k8sva.NewK8SClient("http://127.0.0.1:8200", "your Vault role")
// Check if there is a failure connecting to Vault
if err != nil {
    log.Fatal(err)
}

// Login to Vault
if err := client.Login(); err != nil {
    log.Fatal(err)
}

// Retrieve the Vault Client
// For more information on how to use the Vault Go Client:
// https://github.com/hashicorp/vault
_ = client.GetVaultClient()
```

### Development example
The default K8SClient does not allow for insecure transport, which should be the case in a production environment!
For development you might use an insecure Vault instance.

```go
import (
    k8sva "github.com/martijngerrits/k8s-vault-authentication"
    "log"
    "net/http"
    "time"
)

client, err := k8sva.NewK8SClientWithOptions(
    "http://127.0.0.1:8200",
    "your Vault role",
    "/var/run/secrets/kubernetes.io/serviceaccount/token",
    "/v1/auth/kubernetes/login",
    10 * time.Second,
    &http.Client{},
    true,
)
// Check if there is a failure connecting to Vault
if err != nil {
    log.Fatal(err)
}

// Login to Vault
if err := client.Login(); err != nil {
    log.Fatal(err)
}

// Retrieve the Vault Client
_ = client.GetVaultClient()
```
