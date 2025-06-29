package routers

import (
	"github.com/gin-gonic/gin"

	c "be-weeklytask-ewallet/controllers"
)

func TransactionRouters(r *gin.RouterGroup) {
        r.POST("/transfer", c.TransferHandler)
        r.POST("/topup", c.TopUpHandler)
        r.GET("/history/:id", c.GetHistoryHandler)
}