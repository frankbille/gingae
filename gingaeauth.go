// Gin middleware, handling authentication against the Google App Engine users
// service.
package gingaeauth

import (
	"github.com/gin-gonic/gin"
	"appengine"
	"appengine/user"
)

// Create App Engine authentication middleware.
func GaeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaectx := appengine.NewContext(c.Req)
		u := user.Current(gaectx)
		if u == nil {
			
		} else {
			c.Set("gaeuser", u)
		}
	}
}
