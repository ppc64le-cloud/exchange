package internal

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var _ Backend = &Kubernetes{}

type Kubernetes struct {
	// kubeconfig file with the context to talk to kubernetes cluster
	Kubeconfig string `mapstructure:"kubeconfig"`
	ClientSet  *kubernetes.Clientset
	Config     *restclient.Config
}

func (k *Kubernetes) AdminContext() *kubernetes.Clientset {
	return k.ClientSet
}

func (k *Kubernetes) UserContext(c *gin.Context) (*kubernetes.Clientset, error) {
	session := sessions.Default(c)

	accessToken := session.Get("access_token")
	if accessToken == nil {
		return nil, fmt.Errorf("access_token is not set")
	}
	refreshToken := session.Get("refresh_token")
	if refreshToken == nil {
		return nil, fmt.Errorf("refresh_token is not set")
	}

	fmt.Printf("access_token: %s and refresh_token: %s\n", accessToken, refreshToken)
	userConfig := &restclient.Config{
		Host:            k.Config.Host,
		TLSClientConfig: restclient.TLSClientConfig{Insecure: true},
		AuthProvider: &clientcmdapi.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"client-id":      Config.ClientID,
				"client-secret":  Config.ClientSecret,
				"id-token":       accessToken.(string),
				"idp-issuer-url": Config.IssuerURL,
				"refresh-token":  refreshToken.(string),
			},
		},
	}
	cs, err := kubernetes.NewForConfig(userConfig)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (k *Kubernetes) GetTemplate(c *gin.Context, name string) {
	cs, err := k.UserContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	l, err := cs.CoreV1().Pods("kube-system").List(c.Request.Context(), metav1.ListOptions{})
	spew.Dump(l)
	spew.Dump(err)
	if errors.IsUnauthorized(err) {
		c.JSON(http.StatusUnauthorized, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusAccepted, l)
}
