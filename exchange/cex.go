package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/marklude/logger"
)

type cex struct {
	apiKey    string
	secretKey string
	ws        *websocket.Dialer
	wsURL     string
	logger    logger.Logger
}

type authentication struct {
	E    string `json:"e"`
	Auth *auth  `json:"auth"`
}

type auth struct {
	Key       string `json:"key"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

type Ticker struct {
	E    string   `json:"e"`
	Data []string `json:"data"`
}

type Cex interface {
	GetTicker()
}

func NewCex(apiKey, secretKey, url string) Cex {
	return &cex{
		apiKey:    apiKey,
		secretKey: secretKey,
		ws: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
		wsURL: url,
	}

}

func (c *cex) GetTicker() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	wsocket, _, err := c.ws.Dial(c.wsURL, nil)

	if err != nil {
		c.logger.Error(err)
	}
	timestamp := time.Now().Unix()
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(fmt.Sprintf("%d%s", timestamp, c.apiKey)))
	sha := hex.EncodeToString(h.Sum(nil))

	newAuth := &auth{Key: c.apiKey, Signature: sha, Timestamp: timestamp}
	newAuthentication := &authentication{E: "auth", Auth: newAuth}

	jsonData, err := json.Marshal(newAuthentication)
	if err != nil {
		c.logger.Error(err)
	}

	defer wsocket.Close()

	go func() {
		defer close(done)

		err = wsocket.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			c.logger.Error(err)
		}
		for {
			_, message, err := wsocket.ReadMessage()
			if err != nil {
				c.logger.Error(err)
				return
			}
			c.logger.Info(string(message))
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			jsonTickerData, err := json.Marshal(&Ticker{E: "ticker", Data: []string{"BNB", "USD"}})
			if err != nil {
				c.logger.Error(err)
				return
			}
			err = wsocket.WriteMessage(websocket.TextMessage, jsonTickerData)
			if err != nil {
				c.logger.Error(err)
				return
			}
		case <-interrupt:
			fmt.Println("Interrupted...")
			err := wsocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.Error(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}
}
