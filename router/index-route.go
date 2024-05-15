package router

import (
	"learn/httpserver/controller"
	"learn/httpserver/utils"
	"net/url"

	"github.com/gin-gonic/gin"
	// "learn/httpserver/controller"
	hydra "github.com/ory/hydra-client-go/client"
)

var (
	adminURL, _ = url.Parse("http://localhost:4445")
	hydraClient = hydra.NewHTTPClientWithConfig(nil,
		&hydra.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
)


func IndexRoute(route *gin.Engine) {

	// --------- to accesss this route, you need to pass the hydra token in the header ----------------
	route.GET("/", controller.Get)

	route.POST("/create", controller.Create)

	route.DELETE("/:id", controller.Delete)

	route.PUT("/:id", controller.Update)

	// get refresh token and generate new access token
	route.POST("/refresh-token", utils.VerifyToken(1), controller.RefreshToken)

	// 0 for access token, 1 for refresh token
	route.GET("/employee", utils.VerifyToken(0), controller.GetEmployeeData)

	route.POST("/logout", controller.Logout)

	//service
	route.POST("/new-service", controller.AssignNewServiceToUser)
	// --------- to accesss this route, you need to pass the hydra token in the header ----------------


	// --------- hydra start ----------------

	Hydracontroller := controller.Handler{
		HydraAdmin: hydraClient.Admin,
	}

	route.GET("/oauth2/auth", controller.HydraPublicPortCall)
	route.GET("/login", Hydracontroller.AuthGetLogin)
	route.POST("/login", Hydracontroller.AuthPostLogin)
	route.GET("/consent", Hydracontroller.AuthGetConsent)
	route.POST("/consent", Hydracontroller.AuthPostConsent)
	// call hydra token endpoint
	route.POST("/oauth2/token", Hydracontroller.HydraTokenEndpoint)
	route.POST("/introspect", Hydracontroller.HydraIntroSpectEndpoint)

	
	route.GET("/test", controller.Test)
	// --------- hydra end ----------------



}
