package auth

import (
	tools "mogong/internal/pkg/tools/net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Next() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		refer := c.Request.Header.Get("Referer")
		reqURI := c.Request.RequestURI
		language := c.Request.Header.Get("Accept-Languages")
		// 验证Referer
		if tools.GetReferDomain(refer) != tools.GetReferDomain(reqURI) {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Referer验证失败"})
			c.Abort()
			return
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With, X-CSRF-Token, AccessToken, Token, Accept-Languages ")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
		}
		c.Set("lang", language)
		c.Next()
	}
}
