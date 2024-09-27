package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 异常处理
func Handler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			switch t := r.(type) {
			case *Response[any]:
				log.Printf("panic: %v\n", t.Message)
				c.JSON(t.Code, gin.H{
					"message": t.Message,
				})
			case string:
				log.Printf("panic: internal error")
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": t,
				})
			default:
				log.Printf("panic: internal error")
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "服务器内部异常",
				})
			}
			c.Abort()
		}
	}()

	c.Next()
}
