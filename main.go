package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	myip   string = ""
	secret string = os.Getenv("SECRET")
)

var (
	method = os.Getenv("SS_METHOD")
	passwd = os.Getenv("SS_PASSWD")
	port   = os.Getenv("SS_PORT")
)

const (
	url = "ss://%s@%s:%s?remarks=US_DC6_CN_PROXY"
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
		if token != Md5(secret) {
			ctx.String(200, ctx.GetHeader(gin.PlatformCloudflare))
			return
		}
		mp := Base64(fmt.Sprintf("%s:%s", method, passwd))
		ctx.String(200, Base64(fmt.Sprintf(url, mp, myip, port)))
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
