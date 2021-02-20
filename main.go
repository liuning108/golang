package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"p2pserve/room"
	"p2pserve/server"
	"p2pserve/tencentyun"
	"strings"
)

const (
	sdkappid = 1400483972
	key      = "ba1052e997a037742c3c062f4aa50f37802e21e9c26ba85eda2919693a312042"
)

func main() {

	roomManager := room.NewRoomManager()
	wsServer := server.NewP2PServer(roomManager.HandleNewWebSocket)
	r := gin.Default()
	// 允许使用跨域请求  全局中间件
	r.Use(Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"messagetype": "pong",
		})
	})

	r.GET("/usersig/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		sig, err := tencentyun.GenUserSig(sdkappid, key, userId, 86400*180)
		if err != nil {
			c.JSON(200, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"sdkAppId": sdkappid,
				"userSig":  sig,
			})
		}

	})

	r.GET("/ws", func(c *gin.Context) {
		wsServer.HandleWebSocketRequest(c.Writer, c.Request)
	})
	//r.Run(":8080")
	r.RunTLS(":443", "./2_yunwu.red.crt", "./3_yunwu.red.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
