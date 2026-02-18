package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	email "github.com/xhit/go-simple-mail/v2"
)

type Email struct {
	Subject  string     `json:"asunto,omitempty"`
	Message  string     `json:"contenido,omitempty"`
	From     string     `json:"from,omitempty"`
	To       []string   `json:"to,omitempty"`
	Cc       *[]string  `json:"cc,omitempty"`
	Bcc      *[]string  `json:"bcc,omitempty"`
	Adjuntos *[]Adjunto `json:"adjuntos,omitempty"`
}

type Adjunto struct {
	File []byte `json:"archivo,omitempty"`
	Name string `json:"nombre_archivo,omitempty"`
}

// Metodo para enviar correo
// @ID           Send
// @Summary      metodo para enviar correo.
// @Description  metodo para enviar correo.
// @Tags         Send
// @Accept       json
// @Produce      json
// @Param        request  body  Email  true  "Contenido del correo con adjuntos en base64"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notificacion [post]
func Send(c *gin.Context) {
	var err error
	var body Email
	statusCode := http.StatusOK
	message := "correo enviado exitosamente"

	defer func() {
		errString := ""
		if err != nil {
			errString = err.Error()
			log.Errorf("<<<ERROR: %s --> %v", message, err)
		} else {
			log.Infof("[DEFER Send] correos enviados desde: %v -> to: %v, cc: %v, bcc: %v", body.From, body.To, body.Cc, body.Bcc)
		}
		c.JSON(statusCode, gin.H{
			"message": message,
			"error":   errString,
		})
	}()

	if bindErr := c.ShouldBindJSON(&body); bindErr != nil {
		err = bindErr
		statusCode = http.StatusBadRequest
		message = "error al parsear body del request"
		log.Errorf("<<<ERROR: %s --> %v", message, err)
		return
	}

	if msgAux, errAux := validateRequest(body); errAux != nil {
		statusCode = http.StatusBadRequest
		message = msgAux
		err = errAux
		return
	}

	server := email.NewSMTPClient()
	server.Host = mustGetEnv("SERVER")
	port, portErr := strconv.Atoi(mustGetEnv("SMTP_PORT"))
	if portErr != nil {
		err = portErr
		statusCode = http.StatusInternalServerError
		message = "error al parsear SMTP_PORT"
		log.Errorf("<<<ERROR: %s --> %v", message, err)
		return
	}
	server.Port = port
	server.Encryption = email.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 30 * time.Second
	server.Username = mustGetEnv("SMTP_USERNAME")
	server.Password = mustGetEnv("SMTP_PASSWORD")
	server.Authentication = email.AuthLogin

	timeout, atoiErr := strconv.Atoi(mustGetEnv("TIMEOUT_SECONDS"))
	if atoiErr != nil {
		err = atoiErr
		statusCode = http.StatusInternalServerError
		message = "error al parsear TIMEOUT_SECONDS"
		log.Errorf("<<<ERROR: %s --> %v", message, err)
		return
	}
	server.SendTimeout = time.Duration(timeout) * time.Second

	smtpClient, connectErr := server.Connect()
	if connectErr != nil {
		err = connectErr
		statusCode = http.StatusInternalServerError
		message = "error al conectarse al servidor: " + server.Host
		log.Errorf("<<<ERROR: %s --> %v", message, err)
		return
	}

	msg := email.NewMSG()
	msg.SetFrom(strings.TrimSpace(mustGetEnv("SMTP_USERNAME"))).AddTo(body.To...).SetSubject(body.Subject)

	if body.Cc != nil {
		msg.AddCc((*body.Cc)...)
	}
	if body.Bcc != nil {
		msg.AddBcc((*body.Bcc)...)
	}
	msg.SetBody(email.TextHTML, body.Message)
	if body.Adjuntos != nil {
		for _, att := range *body.Adjuntos {
			name := strings.TrimSpace(att.Name)
			if name == "" {
				continue
			}
			mimeType := filepath.Ext(name)
			msg.AddAttachmentData(att.File, name, mimeType)
		}
	}

	if msg.Error != nil {
		err = msg.Error
		statusCode = http.StatusBadRequest
		message = "error en la validacion antes de enviar correo"
		log.Errorf("<<<ERROR: %s --> %v", message, msg.Error)
		return
	}

	if sendErr := msg.Send(smtpClient); sendErr != nil {
		err = sendErr
		statusCode = http.StatusInternalServerError
		message = "error al enviar correo"
		log.Errorf("<<<ERROR: %s --> %v", message, err)
		return
	}
}

func validateRequest(body Email) (string, error) {
	msjErr := "parametro '%s' es obligatorio"
	if len(body.To) == 0 {
		msj := fmt.Sprintf(msjErr, "to")
		log.Errorf("<<<ERROR: %s", msj)
		return msj, errors.New(msj)
	}
	if strings.TrimSpace(body.Subject) == "" {
		msj := fmt.Sprintf(msjErr, "asunto")
		log.Errorf("<<<ERROR: %s", msj)
		return msj, errors.New(msj)
	}
	if strings.TrimSpace(body.Message) == "" {
		msj := fmt.Sprintf(msjErr, "contenido")
		log.Errorf("<<<ERROR: %s", msj)
		return msj, errors.New(msj)
	}
	return "", nil
}

func mustGetEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(fmt.Sprintf("%s missing!", key))
	}
	return value
}
