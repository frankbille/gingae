# Gin App Engine Middleware [![Build Status](https://travis-ci.org/frankbille/gingae.svg?branch=master)](https://travis-ci.org/frankbille/gingae) [![GoDoc](https://godoc.org/github.com/frankbille/gingae?status.png)](https://godoc.org/github.com/frankbille/gingae)

Gin middlewares providing Google App Engine integrations.

## Provided middlewares

Only checked are provided currently.

- [X] **GAE Context** - Set a variable on the Gin context, containing the GAE Context.
- [X] **GAE User**
  - [X] Set a variable on the Gin context, containing the GAE User, logged in using the standard user authentication.
  - [X] Set a variable on the Gin context, containing the GAE User, logged in using OAuth.
- [ ] **GAE Authentication** - Fail a request with a 401 if user is not authenticated.

## Usage

You always have to include GAE Context, as all the others depend on that.

### GAE Context

```go
package app

import (
	"appengine"
	"github.com/frankbille/gingae"
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.New()
	r.Use(gingae.GaeContext())
	r.GET("/posts", func(c *gin.Context) {
		gaeCtx := c.Get(gingae.Context).(appengine.Context)
		
		// Do stuff which requires the GAE Context
	})
	http.Handle("/", r)
}
```

### GAE User (standard)

```go
package app

import (
	"appengine/user"
	"github.com/frankbille/gingae"
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.New()
	// You always have to include GaeContext, as all the other middlewares depend on it.
	r.Use(gingae.GaeContext())
	r.Use(gingae.GaeUser())
	r.GET("/admin", func(c *gin.Context) {
		gaeUser := c.Get(gingae.User).(user.User)
		
		// Do stuff with the GAE User
	})
	http.Handle("/", r)
}
```

### GAE User (OAuth)

```go
package app

import (
	"appengine/user"
	"github.com/frankbille/gingae"
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.New()
	// You always have to include GaeContext, as all the other middlewares depend on it.
	r.Use(gingae.GaeContext())
	r.Use(gingae.GaeUserOAuth("profile"))
	r.GET("/admin", func(c *gin.Context) {
		if c.Get(gingae.UserOAuthError) != nil {
			// Handle OAuth failures
		}
		
		gaeUser := c.Get(gingae.User).(user.User)
		
		// Do stuff with the GAE User
	})
	http.Handle("/", r)
}
```
