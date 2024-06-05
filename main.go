package main

import (
	"learn/httpserver/router"
	"learn/httpserver/setup"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	setup.LoadEnvVariable()

	r := gin.Default()

	router.IndexRoute(r)

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_BASE_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	r := setupRouter()

	r.Use(CORSMiddleware())

	//Register the standard HandlerFuncs from the net/http/pprof package with the provided gin.Engine.
	pprof.Register(r)

	r.Run()

}
