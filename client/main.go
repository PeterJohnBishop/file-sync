package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateID(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func main() {
	id, _ := GenerateID(8)
	r := gin.Default()

	header := make(http.Header)
	header.Add("X-Client-Id", id)
	u := url.URL{Scheme: "ws", Host: "10.0.0.177:8080", Path: "/ws"}
	go func() {
		const maxAttempts = 3
		const totalRetryDuration = 10 * time.Second
		const retryInterval = totalRetryDuration / (maxAttempts - 1)

		var conn *websocket.Conn
		var err error

		for i := 1; i <= maxAttempts; i++ {
			log.Printf("[WS] Attempt %d: connecting to %s", i, u.String())
			conn, _, err = websocket.DefaultDialer.Dial(u.String(), header)
			if err == nil {
				log.Printf("[WS] Successfully connected on attempt %d", i)
				break
			}

			if i < maxAttempts {
				log.Printf("[WS] Attempt %d failed. Retrying in %v...", i, retryInterval)
				time.Sleep(retryInterval)
			} else {
				log.Fatalf("[WS] Critical Error: Failed to connect after %d attempts: %v", i, err)
			}
		}

		defer conn.Close()

		msg := fmt.Sprintf("%s connected.", id)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Printf("[WS] write_error: %v", err)
			return
		}

		log.Printf("[WS] Entering listen loop...")
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[WS] read_error (connection likely closed): %v", err)
				return
			}
			log.Printf("[WS] received from server: %s (Type: %d)", message, messageType)
		}
	}()

	r.Run(":3000")
}
