package controllers

import (
	m "be-weeklytask-ewallet/models"
	u "be-weeklytask-ewallet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterHandler(ctx *gin.Context) {
	var req m.RegisterRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Email is required",
		})
		return
	}

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Password is required",
		})
		return
	}

	if len(req.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Password must be at least 6 characters",
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Password and Confirm Password did not match",
		})
		return
	}

	if m.IsEmailExist(req.Email) {
		ctx.JSON(http.StatusConflict, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Email already exist",
		})
		return
	}
	userUUID := u.GenerateUUID()
	parseUserUUID, err := uuid.Parse(userUUID)
	err = m.InsertUserToDB(req.Email, req.Password, req.Pin, parseUserUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Internal Server Error: Failed to insert user to database",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, u.Response{
		Status:  "Success",
		Success: true,
		Message: "Register Success",
	})

}

func LoginHandler(ctx *gin.Context) {
	var req m.LoginRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Email is required",
		})
		return
	}

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Password is required",
		})
		return
	}

	if req.Pin == "" {
		ctx.JSON(http.StatusBadRequest, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Pin is required",
		})
		return
	}

	if !m.MatchUserInDatabase(req.Email, req.Password, req.Pin) {
		ctx.JSON(http.StatusUnauthorized, u.Response{
			Status:  "Failed",
			Success: false,
			Message: "Unauthorized. Make sure your email, password and pin is correct",
		})
		return
	}

	ctx.JSON(http.StatusOK, u.Response{
		Status:  "Success",
		Success: true,
		Message: "Login Success",
	})	
}

