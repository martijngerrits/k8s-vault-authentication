package k8s_vault_authentication

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
)

func (c *K8SVaultClient) Login() error {
	log.Printf("Vault login using path %s role %s with jwt [%d bytes]", c.LoginPath, c.Role, len(c.JWT))

	req := c.Client.NewRequest("POST", c.LoginPath)
	if err := req.SetJSONBody(map[string]interface{}{
		"jwt":  c.JWT,
		"role": c.Role,
	}); err != nil {
		return err
	}

	resp, err := c.Client.RawRequest(req)
	if err != nil {
		return fmt.Errorf("could not login to Vault at %s with error: %v", req.URL.String(), err)
	}

	if err := resp.Error(); err != nil {
		return fmt.Errorf("got error from Vault: %v", err)
	}

	var result api.Secret
	if err := resp.DecodeJSON(&result); err != nil {
		return fmt.Errorf("failed to decode JSON response with error: %v", err)
	}

	c.Client.SetToken(result.Auth.ClientToken)
	return nil
}
