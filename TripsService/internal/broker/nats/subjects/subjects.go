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
	TripCreateEvent NatsSubject = "project.<environment>.trips.trip.create"
	TripUpdateEvent NatsSubject = "project.<environment>.trips.trip.update"
	TripDeleteEvent NatsSubject = "project.<environment>.trips.trip.delete"
	TripGetEvent    NatsSubject = "project.<environment>.trips.trip.get"
	TripJoinEvent   NatsSubject = "project.<environment>.trips.trip.join"
	TripCancelEvent NatsSubject = "project.<environment>.trips.trip.cancel"
)
