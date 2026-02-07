package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"

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
	r.GET("/connect", func(c *gin.Context) {
		header := make(http.Header)
		header.Add("X-Client-Id", id)

		u := url.URL{Scheme: "ws", Host: "0.0.0.0:8080", Path: "/ws"}
		log.Printf("connecting to %s", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			log.Fatal("dial_error:", err)
			c.JSON(500, gin.H{"error": "Failed to connect"})
			return
		}

		defer conn.Close()
		msg := fmt.Sprintf("%s connected.", id)
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write_error:", err)
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read_error:", err)
			return
		}

		c.JSON(200, gin.H{"server_response": string(message)})
	})

	r.Run(":3000")
}
