package subjects

import (
	"github.com/hesher116/MyFinalProject/TripsService/internal/config"
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
	TripCreateEvent NatsSubject = "project.<environment>.trips.trip.create"
	TripUpdateEvent NatsSubject = "project.<environment>.trips.trip.update"
	TripDeleteEvent NatsSubject = "project.<environment>.trips.trip.delete"
	TripGetEvent    NatsSubject = "project.<environment>.trips.trip.get"
	TripJoinEvent   NatsSubject = "project.<environment>.trips.trip.join"
	TripCancelEvent NatsSubject = "project.<environment>.trips.trip.cancel"
)
