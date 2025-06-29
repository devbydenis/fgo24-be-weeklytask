package controllers

import (
	"be-weeklytask-ewallet/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransferRequest struct
type TransferRequest struct {
    SenderID    string  `json:"sender_id" binding:"required"`
    ReceiverID  string  `json:"receiver_id" binding:"required"`
    Amount      float64 `json:"amount" binding:"required,gt=0"`
    Description string  `json:"description"`
    Notes       string  `json:"notes"`
}

// TopUpRequest struct
type TopUpRequest struct {
    UserID      string  `json:"user_id" binding:"required"`
    Amount      float64 `json:"amount" binding:"required,gt=0"`
    Description string  `json:"description"`
}

// Transfer - Handle money transfer
func TransferHandler(ctx *gin.Context) {
    var req TransferRequest
    
    // Bind JSON request
    if err := ctx.ShouldBind(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request data",
            "details": err.Error(),
        })
        return
    }

    // Parse UUIDs
    senderID, err := uuid.Parse(req.SenderID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid sender ID",
        })
        return
    }

    receiverID, err := uuid.Parse(req.ReceiverID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid receiver ID",
        })
        return
    }

    // Cek apakah receiver ada
    _, err = models.GetProfileFromDb(receiverID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "Receiver not found",
        })
        return
    }

    // Transfer money
    err = models.Transfer(senderID, receiverID, req.Amount, req.Description, req.Notes)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Transfer failed - insufficient balance",
        })
        return
    }

    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Transfer successful",
        "amount":  req.Amount,
    })
}

// TopUp - Handle balance top up
func  TopUpHandler(ctx *gin.Context) {
    var req TopUpRequest
    
    // Bind JSON request
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request data",
            "details": err.Error(),
        })
        return
    }

    // Parse UUID
    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }

    // Top up balance
    err = models.TopUp(userID, req.Amount, req.Description)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Top up failed",
        })
        return
    }

    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Top up successful",
        "amount":  req.Amount,
    })
}

