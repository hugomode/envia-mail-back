package main

import (
	"fmt"
	"log"
	"os"

	"envia-mail-back/api"
	"envia-mail-back/docs"
	"envia-mail-back/healthcheck"
	"envia-mail-back/middleware"
	"github.com/gin-gonic/gin"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

// @title           Envia Mail Back
// @version         1.0
// @description     API para el envio de correos.
// @termsOfService  http://swagger.io/terms/
// @contact.name    Hugo Carcamo
// @contact.url     https://mesadeservicios.uchile.cl/otrs/index.pl
// @contact.email   hugo.carcamo@uchile.cl
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost:8080"
	}
	docs.SwaggerInfo.Host = host

	r := gin.New()
	prefix := "/api/v1"
	r.Use(
		gin.Recovery(),
		middleware.Cors(),
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{prefix + "/healthcheck"}}),
	)

	v1NoAuth := r.Group(prefix)
	var h healthcheck.Services
	h.Load("Local", "tcp", "127.0.0.1", 8080)
	v1NoAuth.GET("/healthcheck", h.Healthchecks)
	v1NoAuth.POST("/notificacion", api.Send)

	url := gs.URL(fmt.Sprintf("%s/swagger/doc.json", prefix))
	v1NoAuth.GET("/swagger/*any", gs.WrapHandler(sf.Handler, url))

	if err := r.Run(); err != nil {
		log.Fatalf("error levantando servidor: %v", err)
	}
}
