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

