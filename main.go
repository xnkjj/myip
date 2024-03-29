package main

import (
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	myip   string = ""
	secret string = os.Getenv("SECRET")
)

const reg = `^(\d+).(\d+).(\d+).(\d+)$`

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{gin.PlatformCloudflare})

	r.GET("/put/:ip", func(ctx *gin.Context) {
		if secret != ctx.GetHeader("SECRET") {
			ctx.String(403, "Forbidden")
			return
		}
		ip := ctx.Param("ip")
		if regexp.MustCompile(reg).MatchString(ip) {
			myip = ip
		}
		ctx.String(200, myip)
	})

	r.GET("/", func(ctx *gin.Context) {
		if secret == ctx.GetHeader("SECRET") {
			ctx.String(200, myip)
		} else {
			ctx.String(200, ctx.ClientIP())
		}
	})

	r.Run(":80")
}
