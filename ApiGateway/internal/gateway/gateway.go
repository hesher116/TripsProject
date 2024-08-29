package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats/subjects"
	"github.com/nats-io/nats.go"
)

type gatewayModule struct {
	nats *nats.Conn
}

func NewGatewayModule(natsCLI *nats.Conn) *gatewayModule {
	return &gatewayModule{
		nats: natsCLI,
	}
}

func (gm *gatewayModule) InitNatsSubscribers() (err error) {
	_, err = gm.nats.Subscribe(subjects.UserRegEvent.ToString(), gm.HandleRegisterNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserRegEvent: %w", err)
	}

	_, err = gm.nats.Subscribe(subjects.UserAuthEvent.ToString(), gm.HandleAuthorizationNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserAuthEvent: %w", err)
	}

	return
}

//// Обробляє реєстрацію користувача через HTTP-запит
//func (am *gatewayModule) RegisterUserNats(c *gin.Context) {
//	var jsonData map[string]interface{}
//	if err := c.ShouldBindJSON(&jsonData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
//		return
//	}
//
//	response, err := am.nats.Request("UserCreateEvent", encode(jsonData), nats.DefaultTimeout)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration process trouble: " + err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, decode(response.Data))
//}
//
//// Обробляє авторизацію користувача через HTTP-запит
//func (am *gatewayModule) AuthUserNats(c *gin.Context) {
//	var jsonData map[string]interface{}
//	if err := c.ShouldBindJSON(&jsonData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON(AuthUserNats) trouble: " + err.Error()})
//		return
//	}
//
//	response, err := am.nats.Request("UserAuthEvent", encode(jsonData), nats.DefaultTimeout)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization process trouble: " + err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, decode(response.Data))
//}

// HandleRegisterNats Обробляє реєстрацію користувача через NATS повідомлення
func (gm *gatewayModule) HandleRegisterNats(msg *nats.Msg) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &jsonData); err != nil {
		fmt.Println("JSON trouble:", err)
		return
	}

	// Обробка отриманих даних та формування відповіді
	response, err := gm.nats.Request("UserService.Register", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		fmt.Println("registration process trouble:", err)
		return
	}

	_ = msg.Respond(response.Data)
}

// HandleAuthorizationNats Обробляє авторизацію користувача через NATS повідомлення
func (gm *gatewayModule) HandleAuthorizationNats(msg *nats.Msg) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &jsonData); err != nil {
		fmt.Println("JSON trouble:", err)
		return
	}

	// Обробка отриманих даних та формування відповіді
	response, err := gm.nats.Request("UserService.Authorize", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		fmt.Println("authorization process trouble:", err)
		return
	}

	_ = msg.Respond(response.Data)
}

// encode перетворює структуру в байти для передачі через NATS
func encode(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error transform to bytes", err)
		return nil
	}
	return bytes
}

// decode перетворює байти з NATS назад у структуру
func decode(data []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Error transform from bytes", err)
		return nil
	}
	return result
}
