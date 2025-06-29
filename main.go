package main

import (
	"be-weeklytask-ewallet/routers"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	routers.CombineRouters(r)

	godotenv.Load()
	r.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT")))
}