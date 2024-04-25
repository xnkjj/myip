package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	myip   string = ""
	secret string = os.Getenv("SECRET")
	ss     string = os.Getenv("SS_CONFIG")
)

const reg = `^(\d+).(\d+).(\d+).(\d+)$`

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/put/:ip", func(ctx *gin.Context) {
		if secret != ctx.GetHeader("SECRET") {
			ctx.String(403, "Forbidden")
			return
		}
		ip := ctx.Param("ip")
		if regexp.MustCompile(reg).MatchString(ip) {
			myip = ip
		}
		go func(ip string) {
			_ = PutRecord(ctx, ip)
		}(myip)
		ctx.String(200, myip)
	})

	r.GET("/", func(ctx *gin.Context) {
		if secret == ctx.GetHeader("SECRET") {
			ctx.String(200, myip)
		} else {
			ctx.String(200, ctx.GetHeader(gin.PlatformCloudflare))
		}
	})

	r.GET("/listing", func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token != Md5(secret) || len(myip) == 0 {
			ctx.String(200, ctx.GetHeader(gin.PlatformCloudflare))
			return
		}
		res := strings.ReplaceAll(ss, "{ip}", myip)
		ctx.String(200, Base64(res))
	})

	r.Run(":80")
}

func Md5(s string) string {
	ha := md5.New()
	ha.Write([]byte(s))
	return fmt.Sprintf("%x", ha.Sum(nil))
}

func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
