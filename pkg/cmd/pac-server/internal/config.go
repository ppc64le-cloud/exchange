package internal

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Config *config

func init() {
	var err error
	Config, err = NewConfig()
	if err != nil {
		panic(err)
	}
	err = Config.Validate()
	if err != nil {
		panic(err)
	}

	err = Config.Load()
	if err != nil {
		panic(err)
	}
}

type BackendType string

var (
	BackendTypeKubernetes BackendType = "kubernetes"
	BackendTypeManageIQ   BackendType = "manageIQ"
)

type config struct {
	CookieStoreKey string `mapstructure:"store_key"`
	ClientID       string `mapstructure:"client_id"`
	ClientSecret   string `mapstructure:"client_secret"`
	IssuerURL      string `mapstructure:"issuer_url"`
	ClientURL      string `mapstructure:"client_url"`

	// BackendType configuration
	BackendType BackendType `mapstructure:"backend_type"`
	Backend
	Kubernetes `mapstructure:"kubernetes"`
	ManageIQ   `mapstructure:"manageIQ"`
}

func (c *config) Validate() error {
	if c.CookieStoreKey == "" {
		return fmt.Errorf("store_key can't be empty")
	}
	if c.ClientID == "" {
		return fmt.Errorf("client_id can't be empty")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client_secret can't be empty")
	}
	if c.IssuerURL == "" {
		return fmt.Errorf("issuer_url can't be empty")
	}
	if c.ClientURL == "" {
		return fmt.Errorf("client_url can't be empty")
	}

	if c.BackendType != BackendTypeKubernetes && c.BackendType != BackendTypeManageIQ {
		return fmt.Errorf("invalid backend type: \"%s\", only kubernetes and manageIQ are supported as backend_type", c.BackendType)
	}
	return nil
}

func (c *config) Load() error {
	switch c.BackendType {
	case BackendTypeKubernetes:
		if err := c.loadKubernetes(); err != nil {
			return err
		}
	case BackendTypeManageIQ:
		if err := c.loadManageIQ(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid backend type: %s", c.BackendType)
	}
	return nil
}

func (c *config) loadKubernetes() error {
	if c.Kubernetes.Kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			c.Kubernetes.Kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	var err error
	// use the current context in kubeconfig
	c.Kubernetes.Config, err = clientcmd.BuildConfigFromFlags("", c.Kubernetes.Kubeconfig)
	if err != nil {
		return err
	}

	// create the clientset
	c.Kubernetes.ClientSet, err = kubernetes.NewForConfig(c.Kubernetes.Config)
	if err != nil {
		return err
	}

	c.Backend = &c.Kubernetes
	return nil
}
func (c *config) loadManageIQ() error {
	// TODO(mkumatag): Implement me!
	return nil
}

func NewConfig() (*config, error) {
	var c config
	viper.AutomaticEnv()

	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/exchange/")
	viper.AddConfigPath("$HOME/.exchange")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Fatal error config file: %w \n", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("Fatal error while Unmarshal: %w \n", err)
	}
	return &c, nil
}
