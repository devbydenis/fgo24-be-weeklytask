package routers

import (
	c "be-weeklytask-ewallet/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRouters(r *gin.RouterGroup) {
	r.GET("/profile/:id", c.GetProfileHandler)
}