package routers

import "github.com/gin-gonic/gin"

func CombineRouters(r *gin.Engine) {
	AuthRouters(r.Group("/auth"))
}