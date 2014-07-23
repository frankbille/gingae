package gingae

import (
	"appengine"
	"appengine/aetest"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGaeContext(t *testing.T) {
	gaeCtx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer gaeCtx.Close()

	handler := gaeContextFromProvider(func(c *gin.Context) appengine.Context {
		return gaeCtx
	})

	ginCtx := gin.Context{}

	handler(&ginCtx)

	if _, err := ginCtx.Get(Context); err != nil {
		t.Fail()
	}
}

func TestGaeUser(t *testing.T) {
	gaeCtx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer gaeCtx.Close()

	ginCtx := gin.Context{}
	ginCtx.Set(Context, gaeCtx)

	GaeUser()(&ginCtx)

	if _, err := ginCtx.Get(User); err != nil {
		t.Fail()
	}
}

func TestGaeUserOAuth(t *testing.T) {
	gaeCtx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer gaeCtx.Close()

	ginCtx := gin.Context{}
	ginCtx.Set(Context, gaeCtx)

	GaeUserOAuth("")(&ginCtx)

	if _, err := ginCtx.Get(User); err != nil {
		t.Fail()
	}
}
