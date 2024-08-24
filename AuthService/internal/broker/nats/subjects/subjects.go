package subjects

import (
	"os"
	"strings"
)

type NatsSubject string

const (
	envTag = "<environment>"
)

func (sub NatsSubject) ToString() string {
	subj := string(sub)

	env := os.Getenv("ENVIRONMENT")

	return strings.Replace(subj, envTag, strings.ToLower(env), 1)
}

// subjects
const (
	UserRegister      NatsSubject = "project.<environment>.trips.auth.register"
	UserAuthorization NatsSubject = "project.<environment>.trips.auth.authorization"
)
