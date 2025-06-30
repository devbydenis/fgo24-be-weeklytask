package controllers

import (
	"be-weeklytask-ewallet/models"
	"be-weeklytask-ewallet/utils"
	// "fmt"
	"net/http"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadPhotoHandler(ctx *gin.Context){

	userId, err := uuid.Parse(ctx.Request.FormValue("userid"))
	if err != nil  {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success:  false,
			Message: "Bad Request to parse uuid",
		})
		return
	}	

	file, err := ctx.FormFile("file") //nama dari field yang dikirim (kalo diform data -> key nya apa)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success:  false,
			Message: "Failed to upload file",
		})
		return
	}

	if file.Size > 5*1024*1024 {  // kita ambil file size dari meta datanya
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success:  false,
			Message: "File size is too big",
		})
		return
	}
	
	// agar supaya kalo ada yg upload dengan nama yang sama ga ke replace kita lakuin hal berikut
	fileName := uuid.New().String()     // generate random string buat nama filenya
	ext := filepath.Ext(file.Filename)  // ambil extention dari metada file.Filename

	err = ctx.SaveUploadedFile(file, "./uploads/"+fileName+ext) // ini relative ke file main.go (file utama)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success:  false,
			Message: "Failed to save file",
			Errors: err.Error(),
		})
		return
	}

	err = models.ChangeProfileImgDB(fileName+ext, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success:  false,
			Message: "Failed to change profile image",
			Errors: err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, utils.Response{
		Success:  true,
		Message: "Upload file successfully",
	})
}