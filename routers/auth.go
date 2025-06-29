package routers

import (
	c "be-weeklytask-ewallet/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRouters(r *gin.RouterGroup){
	r.POST("/register", c.RegisterHandler)
	r.POST("/login", c.LoginHandler)
}