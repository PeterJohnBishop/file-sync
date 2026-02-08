package server

import (
	"file-sync/server/watcher"
	"file-sync/server/websocket"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ServeGin() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // listen on all network interfaces
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gin.SetMode(gin.ReleaseMode) // remove distracting logging
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			return false
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	hub := websocket.NewHub()
	go hub.Run()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"identity": "File-Sync",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c)
	})

	go watcher.Watch()

	svr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Serving Gin at %s", svr)
	r.Run(svr)
}
