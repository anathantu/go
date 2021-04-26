package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu")
	})

	//API参数，获得路径上的信息
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		action = strings.Trim(action, "/")
		c.String(http.StatusOK, name+" is "+action)
	})

	//Url参数，获得路径后跟随的参数
	r.GET("/user", func(c *gin.Context) {
		//默认值为tutu，/user?name=changchang的话，值为changchang
		name := c.DefaultQuery("name", "tutu")
		c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	})
	r.Run()
}
