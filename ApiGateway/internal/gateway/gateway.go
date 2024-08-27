package gateway

import (
	"fmt"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats/subjects"
	"github.com/nats-io/nats.go"
	"log"
)

type AuthorizationModule struct {
	nats *nats.Conn
}

func NewAuthorizationModule(natsCLI *nats.Conn) *AuthorizationModule {
	return &AuthorizationModule{
		nats: natsCLI,
	}
}

func (am *AuthorizationModule) InitNatsSubscribers() (err error) {
	_, err = am.nats.Subscribe(subjects.UserRegEvent.ToString(), am.RegisterNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserRegEvent: %w", err)

	}

	_, err = am.nats.Subscribe(subjects.UserAuthEvent.ToString(), am.AuthorizationNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserAuthEvent: %w", err)
	}

	return
}

func (am *AuthorizationModule) RegisterNats(m *nats.Msg) {
	log.Printf("RegisterNATS called: %s\n", string(m.Data))

	// Пересилаємо запит до UserService для реєстрації користувача
	response, err := am.nats.Request("UserService.Register", m.Data, nats.DefaultTimeout)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		_ = m.Respond([]byte("Failed to register user: " + err.Error()))
		return
	}

	// Відправляємо відповідь назад до ApiGateway
	_ = m.Respond(response.Data)
}

func (am *AuthorizationModule) AuthorizationNats(m *nats.Msg) {
	log.Printf("AuthorizationNATS called: %s\n", string(m.Data))

	// Пересилаємо запит до UserService для авторизації користувача
	response, err := am.nats.Request("UserService.Authorize", m.Data, nats.DefaultTimeout)
	if err != nil {
		log.Printf("Failed to authorize user: %v", err)
		_ = m.Respond([]byte("Failed to authorize user: " + err.Error()))
		return
	}

	// Відправляємо відповідь назад до ApiGateway
	_ = m.Respond(response.Data)
}
