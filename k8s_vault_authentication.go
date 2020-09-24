package k8s_vault_authentication

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	GetVaultClient() *api.Client
	Login() error
}

type K8SVaultClient struct {
	Client    *api.Client
	Role      string
	JWT       string
	LoginPath string
}

func (c *K8SVaultClient) GetVaultClient() *api.Client {
	return c.Client
}

func lookupJwt(tokenPath string) (string, error) {
	buf, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	s := string(buf)
	return s, nil
}

func NewK8SClient(url string, role string) (Client, error) {
	if strings.HasPrefix(url, "http://") {
		return nil, fmt.Errorf("no insecure transport allowed")
	}

	client, err := api.NewClient(&api.Config{
		Address:    url,
		HttpClient: &http.Client{},
		Timeout:    time.Second * 10,
	})
	if err != nil {
		return nil, err
	}

	jwt, err := lookupJwt("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, err
	}

	return &K8SVaultClient{
		Client:    client,
		Role:      role,
		JWT:       jwt,
		LoginPath: "/v1/auth/kubernetes/login",
	}, nil
}

func NewK8SClientWithOptions(
	url string,
	role string,
	tokenPath string,
	loginPath string,
	timeout time.Duration,
	httpClient *http.Client,
	insecure bool,
) (Client, error) {
	if strings.HasPrefix(url, "http://") && !insecure {
		return nil, fmt.Errorf("no insecure transport allowed, use insecure=true if you are sure")
	}

	client, err := api.NewClient(&api.Config{
		Address:    url,
		HttpClient: httpClient,
		Timeout:    timeout,
	})
	if err != nil {
		return nil, err
	}

	jwt, err := lookupJwt(tokenPath)
	if err != nil {
		return nil, err
	}

	return &K8SVaultClient{
		Client:    client,
		Role:      role,
		JWT:       jwt,
		LoginPath: loginPath,
	}, nil
}
