package exchange

import (
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

type Ticker struct {
	Timestamp string `json:"timestamp"`
	Low       string `json:"low"`
}

type Cex interface {
	getTicker()
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

func (c *cex) getTicker() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	wsocket, _, err := c.ws.Dial(c.wsURL, nil)

	if err != nil {
		c.logger.Error("cex getticker", err)
	}

	defer wsocket.Close()

	go func() {
		defer close(done)
		for {
			_, message, err := wsocket.ReadMessage()
			if err != nil {
				c.logger.Error("read message", err)
				return
			}
			c.logger.Info("received %s :", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := wsocket.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				c.logger.Error("write message error", err)
				return
			}
		case <-interrupt:
			fmt.Println("Interrupted...")
			err := wsocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.Error("interrupted websocket", err)
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
