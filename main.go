package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"p2pserve/server"
	"p2pserve/tencentyun"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	sdkappid = 1400483972
	key      = "ba1052e997a037742c3c062f4aa50f37802e21e9c26ba85eda2919693a312042"
)

func main() {
	allManager := server.NewAllManager()
	roomManager := server.NewRoomManager() //
	roomServer := server.NewWsServer(roomManager.HandleNewWebSocket)
	allServe := server.NewWsServer(allManager.HandleNewWebSocket)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"messagetype": "pong",
		})
	})
	r.GET("/ws/room", func(c *gin.Context) {
		roomServer.HandleWebSocketRequest(c.Writer, c.Request)
	})

	r.GET("/ws", func(c *gin.Context) {
		allServe.HandleWebSocketRequest(c.Writer, c.Request)
	})

	r.GET("/usersig/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		sig, _ := tencentyun.GenUserSig(sdkappid, key, userId, 86400*180)
		c.JSON(200, gin.H{
			"sdkAppId": sdkappid,
			"userSig":  sig,
		})

	})
	r.Run(":80")
	//r.RunTLS(":443","./2_yunwu.red.crt","./3_yunwu.red.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
