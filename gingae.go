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
		return appengine.NewContext(c.Req)
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
		gaeCtx := c.Get(Context).(appengine.Context)
		gaeUser := user.Current(gaeCtx)
		c.Set(User, gaeUser)
	}
}

// Set a variable on the Gin context, containing the GAE User, logged in using OAuth.
func GaeUserOAuth(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := c.Get(Context).(appengine.Context)
		gaeUser, err := user.CurrentOAuth(gaeCtx, scope)
		if err != nil {
			c.Set(UserOAuthError, err)
		} else {
			c.Set(User, gaeUser)
		}
	}
}
