package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)


func main ()  {
	flag.Parse()
	hub := newHub()
	go hub.run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/room/:roomId", func(c *gin.Context){
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId", func(c *gin.Context){
		roomId := c.Param("roomId")
		serveWs(hub, c.Writer, c.Request, roomId)
	})

	err := router.Run("0.0.0.0:8000")
	if err != nil {
		return 
	}
}

