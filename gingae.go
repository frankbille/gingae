// Gin middleware, handling authentication against the Google App Engine users
// service.
package gingae

import (
	"appengine"
	"appengine/user"
	"github.com/gin-gonic/gin"
)

const (
	Context        = "GaeContext"
	User           = "GaeUser"
	UserOAuthError = "GaeUserOAuthError"
)

type gaeContextProvider func(c *gin.Context) appengine.Context

// Set a variable on the Gin context, containing the GAE Context.
func GaeContext() gin.HandlerFunc {
	return gaeContextFromProvider(func(c *gin.Context) appengine.Context {
		return appengine.NewContext(c.Request)
	})
}

func gaeContextFromProvider(gaeContextProvider gaeContextProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := gaeContextProvider(c)
		c.Set(Context, gaeCtx)
	}
}

// Set a variable on the Gin context, containing the GAE User.
func GaeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx, err := c.Get(Context)
		if err != nil {
			panic("Must use the GaeContext middleware before the GaeUser")
		}
		gaeUser := user.Current(gaeCtx.(appengine.Context))
		c.Set(User, gaeUser)
	}
}

// Set a variable on the Gin context, containing the GAE User, logged in using OAuth.
func GaeUserOAuth(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx, err := c.Get(Context)
		if err != nil {
			panic("Must use the GaeContext middleware before the GaeUserOAuth")
		}
		gaeUser, err := user.CurrentOAuth(gaeCtx.(appengine.Context), scope)
		if err != nil {
			c.Set(UserOAuthError, err)
		} else {
			c.Set(User, gaeUser)
		}
	}
}
