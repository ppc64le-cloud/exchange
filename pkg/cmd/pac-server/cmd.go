package pac_server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/internal"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/middleware"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/routes"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = &cobra.Command{
		Use: "Power Access Cloud (PAC) Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
}

func NewPacServerCommand() *cobra.Command {
	return cmd
}

func run() error {
	r := gin.Default()
	r.SetTrustedProxies([]string{})

	// PING api for the health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register the middlewares
	r.Use(middleware.Session())
	r.Use(middleware.OIDC(r))
	r.Use(middleware.K8s())

	// Add routes
	v1router := r.Group("/api/v1")
	routes.AddRoutes(v1router)

	r.GET("/hello", hello)
	r.Run()
	return nil
}

func hello(c *gin.Context) {
	//internal.Config.Kubernetes.AdminContext().CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	cs, err := internal.Config.Kubernetes.UserContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	l, err := cs.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	if errors.IsUnauthorized(err) {
		c.JSON(http.StatusUnauthorized, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusAccepted, l)
}
