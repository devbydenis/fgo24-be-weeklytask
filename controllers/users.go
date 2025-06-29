package controllers

import (
	m "be-weeklytask-ewallet/models"
	"be-weeklytask-ewallet/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



func GetProfileHandler(ctx *gin.Context) {
	fmt.Println("GetProfileHandler hit")
	userId := ctx.Param("id")

	parseToUUID, err := uuid.Parse(userId)
	if err != nil  {
		fmt.Println("GetProfileHandler parse uuid error:", err)
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Failed",
			Success: false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}
	
	data, err := m.GetProfileFromDb(parseToUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Failed",
			Success: false,
			Message: "Internal Server Error",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  "Success",
		Success: true,
		Message: "Success",
		Results: data,
	})
}

func GetBalanceHandler(ctx *gin.Context) {
    // Get user ID from URL parameter
    userIDStr := ctx.Param("id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }

    // Get user
    user, err := m.GetProfileFromDb(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "User not found",
        })
        return
    }

    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "balance": user.Balance,
    })
}