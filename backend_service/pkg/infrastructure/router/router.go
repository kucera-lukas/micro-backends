package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"

	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
)

// Routes of our gin.Engine.
const (
	queryPath      = "/graphql"
	playgroundPath = "/playground"
)

// New creates new gin.Engine.
func New(
	config *env.Config,
	srv http.Handler,
) http.Handler {
	router := gin.Default()

	router.Any(queryPath, func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	var crossOrigin *cors.Cors

	switch config.IsDevelopment() {
	case true:
		playgroundHandler := playground.Handler("GQL Playground", queryPath)

		router.Any(
			playgroundPath,
			func(c *gin.Context) {
				playgroundHandler.ServeHTTP(c.Writer, c.Request)
			},
		)

		crossOrigin = cors.AllowAll()

		gin.SetMode(gin.DebugMode)
	case false:
		crossOrigin = cors.New(cors.Options{ //nolint:exhaustivestruct
			AllowedOrigins: getAllowedOrigins(),
			AllowedHeaders: getAllowedHeaders(),
			Debug:          config.Debug,
		})

		gin.SetMode(gin.ReleaseMode)
	}

	return crossOrigin.Handler(router)
}

func getAllowedOrigins() []string {
	return []string{"*"}
}

func getAllowedHeaders() []string {
	return []string{}
}
