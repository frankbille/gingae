package gingaeauth

import (
	"github.com/gin-gonic/gin"
	"appengine"
	"appengine/user"
)

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
