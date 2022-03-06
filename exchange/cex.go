package exchange

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type cex struct {
	apiKey    string
	secretKey string
	ws        *websocket.Dialer
}

type Ticker struct {
	Timestamp string `json:"timestamp"`
	Low       string `json:"low"`
}

type Cex interface {
	getTicker() string
}

func NewCex(apiKey, secretKey string) Cex {
	return &cex{
		apiKey:    apiKey,
		secretKey: secretKey,
		ws: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}

}

func (c *cex) getTicker() string {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	wsocket, err := c.ws.Dial("", nil)
}
