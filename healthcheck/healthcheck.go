package healthcheck

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: Only TCP healthcheck support for now.
// Add others according to your needs.
type Services struct {
	Services []Service
}

type Service struct {
	Name   string
	Dns    string
	Port   int
	Kind   string
	Status string
}

func (s *Services) Load(name, kind, dns string, port int) {
	if kind == "" {
		kind = "tcp"
	}

	s.Services = append(s.Services, Service{
		Name:   name,
		Dns:    dns,
		Port:   port,
		Kind:   kind,
		Status: "OK",
	})
}

func (s *Services) Healthchecks(c *gin.Context) {
	isHealthy := true

	for x, service := range s.Services {
		if service.Kind == "tcp" && !service.tcpChecker() {
			isHealthy = false
		}
		s.Services[x].Status = service.Status
	}

	if isHealthy {
		c.JSON(http.StatusOK, s.Services)
		return
	}

	c.JSON(http.StatusInternalServerError, s.Services)
}

func (s *Service) tcpChecker() bool {
	s.Status = "NoOK"

	if s.Dns == "" || s.Port <= 0 {
		log.Println("DNS and Port must be provided")
		return false
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.Dns, s.Port), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	s.Status = "OK"
	return true
}
