package gingae

import (
	"appengine"
	"appengine_internal"
	pb "appengine_internal/user"
	"errors"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"testing"
)

type MockGaeContext struct {
	FailOnCall bool
}

func (c *MockGaeContext) Debugf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (c *MockGaeContext) Infof(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (c *MockGaeContext) Warningf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (c *MockGaeContext) Errorf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (c *MockGaeContext) Criticalf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (c *MockGaeContext) Call(service, method string, in, out appengine_internal.ProtoMessage, opts *appengine_internal.CallOptions) error {
	if c.FailOnCall {
		return errors.New("Failing")
	}

	if service == "user" {
		email := "test@example.com"
		authDomain := "example.com"
		userId := "userid"
		isAdmin := false
		out.(*pb.GetOAuthUserResponse).Email = &email
		out.(*pb.GetOAuthUserResponse).AuthDomain = &authDomain
		out.(*pb.GetOAuthUserResponse).UserId = &userId
		out.(*pb.GetOAuthUserResponse).IsAdmin = &isAdmin
	} else {
		log.Printf("Service: %v", service)
	}
	return nil
}

func (c *MockGaeContext) FullyQualifiedAppID() string {
	return ""
}

func (c *MockGaeContext) Request() interface{} {
	header := http.Header{}
	header.Add("X-AppEngine-User-Email", "test@example.com")
	header.Add("X-AppEngine-User-Federated-Identity", "test@example.com")
	return &http.Request{
		Header: header,
	}
}

func TestGaeContext(t *testing.T) {
	Convey("When using the GaeContext middleware", t, func() {
		gaeCtx := &MockGaeContext{}

		handler := gaeContextFromProvider(func(c *gin.Context) appengine.Context {
			return gaeCtx
		})

		ginCtx := gin.Context{}

		handler(&ginCtx)

		Convey("The GAE Context should be set on the Gin Context", func() {
			foundGaeCtx, getErr := ginCtx.Get(Context)

			So(getErr, ShouldBeNil)

			So(foundGaeCtx, ShouldEqual, gaeCtx)
		})
	})
}

func TestGaeUser(t *testing.T) {
	Convey("When using the GaeUser middleware", t, func() {
		gaeCtx := &MockGaeContext{}

		ginCtx := gin.Context{}

		Convey("Panic if GaeContext middleware isn't added before.", func() {
			userFunc := func() {
				GaeUser()(&ginCtx)
			}

			So(userFunc, ShouldPanic)
		})

		Convey("The GAE User should be set on the Gin Context", func() {
			ginCtx.Set(Context, gaeCtx)

			GaeUser()(&ginCtx)

			user, getErr := ginCtx.Get(User)

			So(getErr, ShouldBeNil)

			So(user, ShouldNotBeNil)
		})
	})
}

func TestGaeUserOAuth(t *testing.T) {
	Convey("When using the GaeUserOAuth middleware", t, func() {
		gaeCtx := &MockGaeContext{}

		ginCtx := gin.Context{}

		Convey("Panic if GaeContext middleware isn't added before.", func() {
			userFunc := func() {
				GaeUserOAuth("")(&ginCtx)
			}

			So(userFunc, ShouldPanic)
		})

		Convey("The GAE User should be set on the Gin Context", func() {
			ginCtx.Set(Context, gaeCtx)

			GaeUserOAuth("")(&ginCtx)

			user, getErr := ginCtx.Get(User)

			So(getErr, ShouldBeNil)

			So(user, ShouldNotBeNil)
		})

		Convey("If user lookup fails", func() {
			gaeCtx = &MockGaeContext{
				FailOnCall: true,
			}

			ginCtx = gin.Context{}
			ginCtx.Set(Context, gaeCtx)

			GaeUserOAuth("")(&ginCtx)

			Convey("User should not be set on Gin Context, but User error should", func() {
				user, getErr := ginCtx.Get(User)

				So(getErr, ShouldNotBeNil)

				So(user, ShouldBeNil)
			})
		})
	})
}
