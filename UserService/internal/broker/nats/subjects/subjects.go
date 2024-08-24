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
	UserCreateEvent NatsSubject = "project.<environment>.trips.user.create"
	UserEditEvent   NatsSubject = "project.<environment>.trips.user.edit"
)
