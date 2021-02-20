package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"p2pserve/room"
	"p2pserve/server"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	roomManager := room.NewRoomManager() //
	wsServer := server.NewP2PServer(roomManager.HandleNewWebSocket)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"messagetype": "pong",
		})
	})
	r.GET("/ws", func(c *gin.Context) {
		wsServer.HandleWebSocketRequest(c.Writer, c.Request)
	})
	r.Run(":80")
	//r.RunTLS(":443","./2_yunwu.red.crt","./3_yunwu.red.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
