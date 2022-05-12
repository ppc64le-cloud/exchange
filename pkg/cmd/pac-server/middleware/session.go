package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/internal"
)

func Session() gin.HandlerFunc {
	// TODO(mkumatag): memstore is meant for only test and development, need to explore other options for the production purpose
	store := memstore.NewStore([]byte(internal.Config.CookieStoreKey))
	return sessions.Sessions("mysession", store)
}
