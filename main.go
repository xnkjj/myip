package main

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

var myip string = ""

const reg = `^(\d+).(\d+).(\d+).(\d+)$`

func main() {
	r := gin.Default()
	gin.DisableConsoleColor()

	r.GET("/:ip", func(ctx *gin.Context) {
		ip := ctx.Param("ip")
		if regexp.MustCompile(reg).MatchString(ip) {
			myip = ip
			ctx.String(200, myip)
		} else {
			ctx.String(200, myip)
		}
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, myip)
	})

	r.Run(":80")
}
