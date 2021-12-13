package middleware

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/app/v2/models"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("sub")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		ok = enforcer.Enforce(user.(models.User), ctx.Request.URL.Path, ctx.Request.Method)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ou are not allowed to access this resource"})
			return
		}
		ctx.Next()
	}
}
