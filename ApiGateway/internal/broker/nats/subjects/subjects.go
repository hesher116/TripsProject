package subjects

import (
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/config"
	"log"
	"strings"
)

type NatsSubject string

const (
	envTag = "<environment>"
)

func (sub NatsSubject) ToString() string {
	cfg := config.LoadConfig()

	subj := string(sub)

	env := cfg.Envi
	if env == "" {
		log.Fatalf("ENVIRONMENT variable is not set")
	}

	return strings.Replace(subj, envTag, strings.ToLower(env), 1)
}

// subjects
const (
	UserRegEvent  NatsSubject = "project.<environment>.trips.api-user.register"
	UserAuthEvent NatsSubject = "project.<environment>.trips.api-user.auth"
)
